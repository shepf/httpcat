package mcp

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"httpcat/internal/common"
	"httpcat/internal/common/utils"
	"httpcat/internal/common/ylog"
	"httpcat/internal/models"
	"httpcat/internal/storage/auth"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	uuid "github.com/satori/go.uuid"
)

// MCPServer MCP Server 实例
type MCPServer struct {
	sseServer *mcpserver.SSEServer
}

// deleteConfirmTokens 存储删除确认 Token（带过期时间）
var (
	deleteConfirmTokens = make(map[string]deleteConfirmInfo)
	deleteTokenMutex    sync.RWMutex
)

type deleteConfirmInfo struct {
	Filename  string
	ExpiresAt time.Time
}

// NewMCPServer 创建新的 MCP Server
func NewMCPServer() *MCPServer {
	// 创建 MCP Server
	s := mcpserver.NewMCPServer(
		"HttpCat",
		common.Version,
		mcpserver.WithResourceCapabilities(true, true),
		mcpserver.WithLogging(),
	)

	// 注册 Tools
	registerTools(s)

	// 注册 Resources
	registerResources(s)

	// 创建 SSE Server
	sseServer := mcpserver.NewSSEServer(
		s,
		mcpserver.WithBaseURL(""),
		mcpserver.WithStaticBasePath("/mcp"),
	)

	// 启动过期 Token 清理协程
	go cleanupExpiredTokens()

	return &MCPServer{sseServer: sseServer}
}

// GetHandler 获取 Gin 路由处理器（带认证）
func (m *MCPServer) GetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 验证 MCP Auth Token（如果配置了）
		if common.McpAuthToken != "" {
			authHeader := c.GetHeader("Authorization")
			expectedToken := "Bearer " + common.McpAuthToken

			if authHeader != expectedToken {
				ylog.Warnf("MCP", "Unauthorized access attempt from %s", c.ClientIP())
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":   "Unauthorized",
					"message": "Invalid or missing Authorization header. Use 'Bearer <token>' format.",
				})
				return
			}
		}

		m.sseServer.ServeHTTP(c.Writer, c.Request)
	}
}

// ServeHTTP 实现 http.Handler 接口
func (m *MCPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.sseServer.ServeHTTP(w, r)
}

// cleanupExpiredTokens 定期清理过期的删除确认 Token
func cleanupExpiredTokens() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		deleteTokenMutex.Lock()
		now := time.Now()
		for token, info := range deleteConfirmTokens {
			if now.After(info.ExpiresAt) {
				delete(deleteConfirmTokens, token)
			}
		}
		deleteTokenMutex.Unlock()
	}
}

// generateConfirmToken 生成删除确认 Token
func generateConfirmToken(filename string) string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	token := hex.EncodeToString(bytes)

	deleteTokenMutex.Lock()
	deleteConfirmTokens[token] = deleteConfirmInfo{
		Filename:  filename,
		ExpiresAt: time.Now().Add(5 * time.Minute), // 5分钟有效期
	}
	deleteTokenMutex.Unlock()

	return token
}

// verifyConfirmToken 验证删除确认 Token
func verifyConfirmToken(token, filename string) bool {
	deleteTokenMutex.RLock()
	info, exists := deleteConfirmTokens[token]
	deleteTokenMutex.RUnlock()

	if !exists {
		return false
	}

	if time.Now().After(info.ExpiresAt) {
		// Token 已过期，删除它
		deleteTokenMutex.Lock()
		delete(deleteConfirmTokens, token)
		deleteTokenMutex.Unlock()
		return false
	}

	if info.Filename != filename {
		return false
	}

	// 验证成功后删除 Token（一次性使用）
	deleteTokenMutex.Lock()
	delete(deleteConfirmTokens, token)
	deleteTokenMutex.Unlock()

	return true
}

// validateAndResolvePath 安全地验证和解析路径（防止路径遍历和符号链接攻击）
func validateAndResolvePath(baseDir, userPath string) (string, error) {
	return common.ResolvePathWithinBase(baseDir, userPath)
}

// registerTools 注册所有 Tools
func registerTools(s *mcpserver.MCPServer) {
	// 1. 文件列表查询
	s.AddTool(
		mcp.NewTool("list_files",
			mcp.WithDescription("获取上传目录中的文件列表，支持按时间排序"),
			mcp.WithString("dir", mcp.Description("子目录路径，默认为空表示根目录")),
			mcp.WithNumber("limit", mcp.Description("返回文件数量限制，默认50")),
		),
		handleListFiles,
	)

	// 2. 文件信息查询
	s.AddTool(
		mcp.NewTool("get_file_info",
			mcp.WithDescription("获取指定文件的详细信息，包括大小、修改时间、MD5等"),
			mcp.WithString("filename", mcp.Required(), mcp.Description("文件名")),
		),
		handleGetFileInfo,
	)

	// 3. 上传历史查询
	s.AddTool(
		mcp.NewTool("get_upload_history",
			mcp.WithDescription("获取文件上传历史记录，支持分页和筛选"),
			mcp.WithNumber("page", mcp.Description("页码，默认1")),
			mcp.WithNumber("page_size", mcp.Description("每页数量，默认20")),
			mcp.WithString("filename", mcp.Description("按文件名筛选")),
			mcp.WithString("file_md5", mcp.Description("按文件MD5筛选")),
		),
		handleGetUploadHistory,
	)

	// 4. 磁盘使用情况
	s.AddTool(
		mcp.NewTool("get_disk_usage",
			mcp.WithDescription("获取上传目录的磁盘使用情况"),
		),
		handleGetDiskUsage,
	)

	// 5. 请求删除文件（两步确认）
	s.AddTool(
		mcp.NewTool("request_delete_file",
			mcp.WithDescription("请求删除文件（第一步：获取确认Token）。返回一个5分钟有效的确认Token，需要调用 confirm_delete_file 完成删除"),
			mcp.WithString("filename", mcp.Required(), mcp.Description("要删除的文件名")),
		),
		handleRequestDeleteFile,
	)

	// 6. 确认删除文件（两步确认）
	s.AddTool(
		mcp.NewTool("confirm_delete_file",
			mcp.WithDescription("确认删除文件（第二步：使用确认Token执行删除）"),
			mcp.WithString("filename", mcp.Required(), mcp.Description("要删除的文件名")),
			mcp.WithString("confirm_token", mcp.Required(), mcp.Description("request_delete_file 返回的确认Token")),
		),
		handleConfirmDeleteFile,
	)

	// 7. 上传文件
	s.AddTool(
		mcp.NewTool("upload_file",
			mcp.WithDescription("上传文件到 HttpCat（需要 UploadToken）"),
			mcp.WithString("filename", mcp.Required(), mcp.Description("文件名，包含扩展名")),
			mcp.WithString("content_base64", mcp.Required(), mcp.Description("文件内容的 Base64 编码")),
			mcp.WithString("upload_token", mcp.Required(), mcp.Description("上传 Token，格式: appkey:signature:policy")),
		),
		handleUploadFile,
	)

	// 10. 上传图片（写入图片管理，生成缩略图，可在 listThumbImages 中查看）
	s.AddTool(
		mcp.NewTool("upload_image",
			mcp.WithDescription("上传图片到 HttpCat 图片管理（会生成缩略图，可在图片管理页面和 listThumbImages 接口中查看）"),
			mcp.WithString("filename", mcp.Required(), mcp.Description("图片文件名，包含扩展名（如 photo.png）")),
			mcp.WithString("content_base64", mcp.Required(), mcp.Description("图片内容的 Base64 编码")),
			mcp.WithString("upload_token", mcp.Required(), mcp.Description("上传 Token，格式: appkey:signature:policy")),
		),
		handleUploadImage,
	)

	// 注：create_upload_token 功能已从 MCP 移除
	// Token 生成请使用 HTTP API: POST /api/v1/user/createUploadToken (需要登录)

	// 8. 获取统计信息
	s.AddTool(
		mcp.NewTool("get_statistics",
			mcp.WithDescription("获取上传和下载统计信息"),
			mcp.WithString("type", mcp.Description("统计类型：upload/download/all，默认all")),
		),
		handleGetStatistics,
	)

	// 9. 验证文件 MD5
	s.AddTool(
		mcp.NewTool("verify_file_md5",
			mcp.WithDescription("验证文件的 MD5 值是否匹配"),
			mcp.WithString("filename", mcp.Required(), mcp.Description("文件名")),
			mcp.WithString("expected_md5", mcp.Required(), mcp.Description("期望的 MD5 值")),
		),
		handleVerifyFileMD5,
	)

	// 11. 创建文件夹（v0.5.0+）
	s.AddTool(
		mcp.NewTool("create_folder",
			mcp.WithDescription("在上传目录中创建新文件夹"),
			mcp.WithString("name", mcp.Required(), mcp.Description("文件夹名称")),
			mcp.WithString("dir", mcp.Description("父目录路径，默认为空表示根目录")),
		),
		handleCreateFolder,
	)

	// 12. 重命名文件或文件夹（v0.5.0+）
	s.AddTool(
		mcp.NewTool("rename_file",
			mcp.WithDescription("重命名文件或文件夹"),
			mcp.WithString("old_name", mcp.Required(), mcp.Description("原文件/文件夹名称")),
			mcp.WithString("new_name", mcp.Required(), mcp.Description("新文件/文件夹名称")),
			mcp.WithString("dir", mcp.Description("所在目录路径，默认为空表示根目录")),
		),
		handleRenameFile,
	)

	// 13. 批量删除文件/文件夹（v0.5.0+）
	s.AddTool(
		mcp.NewTool("batch_delete_files",
			mcp.WithDescription("批量删除文件或空文件夹（文件夹必须为空才能删除）"),
			mcp.WithString("files", mcp.Required(), mcp.Description("要删除的文件名列表，JSON 数组格式，如 [\"file1.txt\",\"file2.txt\"]")),
			mcp.WithString("dir", mcp.Description("所在目录路径，默认为空表示根目录")),
		),
		handleBatchDeleteFiles,
	)

	// 14. 文件总览统计（v0.5.0+）
	s.AddTool(
		mcp.NewTool("get_file_overview",
			mcp.WithDescription("获取文件总览统计信息：文件总数、目录数、总大小"),
		),
		handleGetFileOverview,
	)

	// 15. 下载历史日志（v0.5.0+）
	s.AddTool(
		mcp.NewTool("get_download_history",
			mcp.WithDescription("获取文件下载历史记录，支持分页和筛选"),
			mcp.WithNumber("page", mcp.Description("页码，默认1")),
			mcp.WithNumber("page_size", mcp.Description("每页数量，默认20")),
			mcp.WithString("filename", mcp.Description("按文件名筛选")),
			mcp.WithString("file_md5", mcp.Description("按文件MD5筛选")),
			mcp.WithString("ip", mcp.Description("按IP地址筛选")),
		),
		handleGetDownloadHistory,
	)
}

// registerResources 注册所有 Resources
func registerResources(s *mcpserver.MCPServer) {
	// 1. 文件列表资源
	s.AddResource(
		mcp.NewResource(
			"filelist://current",
			"Current Files",
			mcp.WithResourceDescription("当前上传目录中的文件列表"),
			mcp.WithMIMEType("application/json"),
		),
		handleFileListResource,
	)

	// 2. 磁盘使用资源
	s.AddResource(
		mcp.NewResource(
			"disk://usage",
			"Disk Usage",
			mcp.WithResourceDescription("磁盘使用情况信息"),
			mcp.WithMIMEType("application/json"),
		),
		handleDiskUsageResource,
	)

	// 3. 系统信息资源
	s.AddResource(
		mcp.NewResource(
			"system://info",
			"System Info",
			mcp.WithResourceDescription("HttpCat 系统信息"),
			mcp.WithMIMEType("application/json"),
		),
		handleSystemInfoResource,
	)
}

// ==================== Tool Handlers ====================

// handleListFiles 处理文件列表查询
func handleListFiles(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	dir := ""
	if dirVal, ok := args["dir"].(string); ok {
		dir = dirVal
	}

	limit := 50
	if limitVal, ok := args["limit"].(float64); ok {
		limit = int(limitVal)
	}

	// 使用安全的路径验证
	var dirPath string
	var err error
	if dir == "" {
		dirPath = common.GetDownloadDir()
	} else {
		dirPath, err = validateAndResolvePath(common.GetDownloadDir(), dir)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid directory path: %v", err)), nil
		}
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return mcp.NewToolResultError("Directory does not exist"), nil
		}
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read directory: %v", err)), nil
	}

	type fileInfo struct {
		Name         string `json:"name"`
		Size         string `json:"size"`
		SizeBytes    int64  `json:"size_bytes"`
		LastModified string `json:"last_modified"`
	}

	// 收集文件信息（过滤掉目录）并按时间倒序排序
	type sortableFile struct {
		info    fileInfo
		modTime time.Time
	}
	var sortableFiles []sortableFile
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		info, err := file.Info()
		if err != nil {
			continue
		}
		sortableFiles = append(sortableFiles, sortableFile{
			info: fileInfo{
				Name:         file.Name(),
				Size:         utils.FormatSize(info.Size()),
				SizeBytes:    info.Size(),
				LastModified: info.ModTime().Format("2006-01-02 15:04:05"),
			},
			modTime: info.ModTime(),
		})
	}

	// 按修改时间倒序排序
	for i := 0; i < len(sortableFiles); i++ {
		for j := i + 1; j < len(sortableFiles); j++ {
			if sortableFiles[j].modTime.After(sortableFiles[i].modTime) {
				sortableFiles[i], sortableFiles[j] = sortableFiles[j], sortableFiles[i]
			}
		}
	}

	var fileList []fileInfo
	for _, sf := range sortableFiles {
		fileList = append(fileList, sf.info)
	}

	// 限制返回数量
	if len(fileList) > limit {
		fileList = fileList[:limit]
	}

	result, _ := json.MarshalIndent(fileList, "", "  ")
	return mcp.NewToolResultText(string(result)), nil
}

// handleGetFileInfo 处理文件信息查询
func handleGetFileInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	filename, ok := args["filename"].(string)
	if !ok || filename == "" {
		return mcp.NewToolResultError("filename is required"), nil
	}

	// 使用安全的路径验证
	filePath, err := validateAndResolvePath(common.GetDownloadDir(), filename)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid filename: %v", err)), nil
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return mcp.NewToolResultError("File not found"), nil
		}
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get file info: %v", err)), nil
	}

	// 计算 MD5
	md5Hash, err := utils.CalculateMD5(filePath)
	if err != nil {
		md5Hash = "N/A"
	}

	type FileDetail struct {
		Name         string `json:"name"`
		Size         string `json:"size"`
		SizeBytes    int64  `json:"size_bytes"`
		LastModified string `json:"last_modified"`
		MD5          string `json:"md5"`
		IsDir        bool   `json:"is_dir"`
	}

	detail := FileDetail{
		Name:         fileInfo.Name(),
		Size:         utils.FormatSize(fileInfo.Size()),
		SizeBytes:    fileInfo.Size(),
		LastModified: fileInfo.ModTime().Format("2006-01-02 15:04:05"),
		MD5:          md5Hash,
		IsDir:        fileInfo.IsDir(),
	}

	result, _ := json.MarshalIndent(detail, "", "  ")
	return mcp.NewToolResultText(string(result)), nil
}

// handleGetUploadHistory 处理上传历史查询
func handleGetUploadHistory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !common.EnableSqlite {
		return mcp.NewToolResultError("SQLite is not enabled"), nil
	}

	args := request.GetArguments()

	page := 1
	pageSize := 20
	filename := ""
	fileMD5 := ""

	if v, ok := args["page"].(float64); ok {
		page = int(v)
	}
	if v, ok := args["page_size"].(float64); ok {
		pageSize = int(v)
	}
	if v, ok := args["filename"].(string); ok {
		filename = v
	}
	if v, ok := args["file_md5"].(string); ok {
		fileMD5 = v
	}

	logs, total, err := queryUploadHistory(page, pageSize, filename, fileMD5, "")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to query history: %v", err)), nil
	}

	response := map[string]interface{}{
		"list":      logs,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	}

	result, _ := json.MarshalIndent(response, "", "  ")
	return mcp.NewToolResultText(string(result)), nil
}

// queryUploadHistory 查询上传历史
func queryUploadHistory(page, pageSize int, filename, fileMD5, ip string) ([]map[string]interface{}, int, error) {
	// 使用 gorm 查询数据库
	db, err := common.GetDB()
	if err != nil {
		return nil, 0, err
	}

	type UploadLog struct {
		ID         int    `gorm:"column:id"`
		IP         string `gorm:"column:ip"`
		AppKey     string `gorm:"column:appkey"`
		UploadTime string `gorm:"column:upload_time"`
		Filename   string `gorm:"column:filename"`
		FileSize   string `gorm:"column:file_size"`
		FileMD5    string `gorm:"column:file_md5"`
	}

	var logs []UploadLog
	offset := (page - 1) * pageSize

	query := db.Table("t_upload_log").Offset(offset).Limit(pageSize).Order("upload_time DESC")
	countQuery := db.Table("t_upload_log")

	if filename != "" {
		query = query.Where("filename LIKE ?", "%"+filename+"%")
		countQuery = countQuery.Where("filename LIKE ?", "%"+filename+"%")
	}
	if fileMD5 != "" {
		query = query.Where("file_md5 = ?", fileMD5)
		countQuery = countQuery.Where("file_md5 = ?", fileMD5)
	}
	if ip != "" {
		query = query.Where("ip = ?", ip)
		countQuery = countQuery.Where("ip = ?", ip)
	}

	if err := query.Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var result []map[string]interface{}
	for _, log := range logs {
		result = append(result, map[string]interface{}{
			"id":          log.ID,
			"ip":          log.IP,
			"appkey":      log.AppKey,
			"upload_time": log.UploadTime,
			"filename":    log.Filename,
			"file_size":   log.FileSize,
			"file_md5":    log.FileMD5,
		})
	}

	return result, int(total), nil
}

// handleGetDiskUsage 处理磁盘使用查询
func handleGetDiskUsage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	usage, err := getDiskUsage(common.GetUploadDir())
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get disk usage: %v", err)), nil
	}

	result, _ := json.MarshalIndent(usage, "", "  ")
	return mcp.NewToolResultText(string(result)), nil
}

// handleRequestDeleteFile 处理删除文件请求（第一步：生成确认Token）
func handleRequestDeleteFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	filename, ok := args["filename"].(string)
	if !ok || filename == "" {
		return mcp.NewToolResultError("filename is required"), nil
	}

	// 使用安全的路径验证
	filePath, err := validateAndResolvePath(common.GetDownloadDir(), filename)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid filename: %v", err)), nil
	}

	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return mcp.NewToolResultError("File not found"), nil
	}
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to check file: %v", err)), nil
	}

	// 生成确认 Token
	confirmToken := generateConfirmToken(filename)

	type DeleteRequest struct {
		Filename     string `json:"filename"`
		Size         string `json:"size"`
		ConfirmToken string `json:"confirm_token"`
		ExpiresIn    string `json:"expires_in"`
		NextStep     string `json:"next_step"`
	}

	response := DeleteRequest{
		Filename:     filename,
		Size:         utils.FormatSize(fileInfo.Size()),
		ConfirmToken: confirmToken,
		ExpiresIn:    "5 minutes",
		NextStep:     "Call confirm_delete_file with filename and confirm_token to complete deletion",
	}

	ylog.Infof("MCP", "Delete requested for file: %s, token: %s", filename, confirmToken[:8]+"...")

	result, _ := json.MarshalIndent(response, "", "  ")
	return mcp.NewToolResultText(string(result)), nil
}

// handleConfirmDeleteFile 处理删除文件确认（第二步：验证Token并执行删除）
func handleConfirmDeleteFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	filename, ok := args["filename"].(string)
	if !ok || filename == "" {
		return mcp.NewToolResultError("filename is required"), nil
	}

	confirmToken, ok := args["confirm_token"].(string)
	if !ok || confirmToken == "" {
		return mcp.NewToolResultError("confirm_token is required"), nil
	}

	// 验证确认 Token
	if !verifyConfirmToken(confirmToken, filename) {
		return mcp.NewToolResultError("Invalid or expired confirm_token. Please call request_delete_file again to get a new token."), nil
	}

	// 使用安全的路径验证
	filePath, err := validateAndResolvePath(common.GetDownloadDir(), filename)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid filename: %v", err)), nil
	}

	// 再次检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return mcp.NewToolResultError("File not found"), nil
	}

	// 执行删除
	if err := os.Remove(filePath); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete file: %v", err)), nil
	}

	ylog.Infof("MCP", "File deleted via MCP (confirmed): %s", filename)

	return mcp.NewToolResultText(fmt.Sprintf("File '%s' deleted successfully", filename)), nil
}

// handleUploadFile 处理文件上传
func handleUploadFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !common.FileUploadEnable {
		return mcp.NewToolResultError("File upload is disabled"), nil
	}

	args := request.GetArguments()

	filename, ok := args["filename"].(string)
	if !ok || filename == "" {
		return mcp.NewToolResultError("filename is required"), nil
	}

	contentBase64, ok := args["content_base64"].(string)
	if !ok || contentBase64 == "" {
		return mcp.NewToolResultError("content_base64 is required"), nil
	}

	uploadToken, ok := args["upload_token"].(string)
	if !ok || uploadToken == "" {
		return mcp.NewToolResultError("upload_token is required"), nil
	}

	// Base64 解码
	content, err := base64.StdEncoding.DecodeString(contentBase64)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid base64 content: %v", err)), nil
	}

	filename, err = common.NormalizeSafeFileName(filename)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid filename: %v", err)), nil
	}
	if strings.HasPrefix(filename, ".") {
		return mcp.NewToolResultError("Invalid filename: hidden files not allowed"), nil
	}

	// 验证 UploadToken
	if common.EnableUploadToken {
		parts := strings.Split(uploadToken, ":")
		if len(parts) != 3 {
			return mcp.NewToolResultError("Invalid UploadToken format"), nil
		}

		appkey := parts[0]

		// 从数据库查询 appsecret
		db, err := common.GetDB()
		if err != nil {
			return mcp.NewToolResultError("Failed to verify token"), nil
		}

		type TokenItem struct {
			Appsecret string `gorm:"column:app_secret"`
			State     string `gorm:"column:state"`
		}

		var tokenItem TokenItem
		result := db.Table("t_upload_token").Where("appkey = ?", appkey).First(&tokenItem)
		if result.Error != nil {
			return mcp.NewToolResultError("Invalid appkey"), nil
		}

		if tokenItem.State == "closed" {
			return mcp.NewToolResultError("Appkey is disabled"), nil
		}

		// 验证 Token
		mac := auth.New(appkey, tokenItem.Appsecret)
		if !mac.VerifyUploadToken(uploadToken) {
			return mcp.NewToolResultError("Invalid UploadToken"), nil
		}
	}

	// 确保上传目录存在
	if err := os.MkdirAll(common.GetUploadDir(), 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create upload directory: %v", err)), nil
	}

	filePath, err := validateAndResolvePath(common.GetUploadDir(), filename)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid filename: %v", err)), nil
	}

	// 写入文件
	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to write file: %v", err)), nil
	}

	// 获取文件信息
	fileInfo, _ := os.Stat(filePath)
	fileSize := utils.FormatSize(fileInfo.Size())
	fileMD5, _ := utils.CalculateMD5(filePath)

	// 记录到数据库（如果启用）
	if common.EnableSqlite {
		go func() {
			db, err := common.GetDB()
			if err != nil {
				return
			}
			db.Table("t_upload_log").Create(map[string]interface{}{
				"ip":          "MCP",
				"appkey":      "",
				"upload_time": time.Now().Format("2006-01-02 15:04:05"),
				"filename":    filename,
				"file_size":   fileSize,
				"file_md5":    fileMD5,
			})
		}()
	}

	ylog.Infof("MCP", "File uploaded via MCP: %s (%s)", filename, fileSize)

	// 构建返回结果
	type UploadResult struct {
		Filename string `json:"filename"`
		Size     string `json:"size"`
		MD5      string `json:"md5"`
		Path     string `json:"path"`
	}

	uploadResult := UploadResult{
		Filename: filename,
		Size:     fileSize,
		MD5:      fileMD5,
		Path:     filePath,
	}

	data, _ := json.MarshalIndent(uploadResult, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

// handleUploadImage 处理图片上传（写入图片管理，生成缩略图）
func handleUploadImage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !common.FileUploadEnable {
		return mcp.NewToolResultError("File upload is disabled"), nil
	}

	args := request.GetArguments()

	filename, ok := args["filename"].(string)
	if !ok || filename == "" {
		return mcp.NewToolResultError("filename is required"), nil
	}

	contentBase64, ok := args["content_base64"].(string)
	if !ok || contentBase64 == "" {
		return mcp.NewToolResultError("content_base64 is required"), nil
	}

	uploadToken, ok := args["upload_token"].(string)
	if !ok || uploadToken == "" {
		return mcp.NewToolResultError("upload_token is required"), nil
	}

	// Base64 解码
	content, err := base64.StdEncoding.DecodeString(contentBase64)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid base64 content: %v", err)), nil
	}

	filename, err = common.NormalizeSafeFileName(filename)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid filename: %v", err)), nil
	}
	if strings.HasPrefix(filename, ".") {
		return mcp.NewToolResultError("Invalid filename: hidden files not allowed"), nil
	}

	// 验证 UploadToken
	if common.EnableUploadToken {
		parts := strings.Split(uploadToken, ":")
		if len(parts) != 3 {
			return mcp.NewToolResultError("Invalid UploadToken format"), nil
		}

		appkey := parts[0]

		db, err := common.GetDB()
		if err != nil {
			return mcp.NewToolResultError("Failed to verify token"), nil
		}

		type TokenItem struct {
			Appsecret string `gorm:"column:app_secret"`
			State     string `gorm:"column:state"`
		}

		var tokenItem TokenItem
		result := db.Table("t_upload_token").Where("appkey = ?", appkey).First(&tokenItem)
		if result.Error != nil {
			return mcp.NewToolResultError("Invalid appkey"), nil
		}

		if tokenItem.State == "closed" {
			return mcp.NewToolResultError("Appkey is disabled"), nil
		}

		mac := auth.New(appkey, tokenItem.Appsecret)
		if !mac.VerifyUploadToken(uploadToken) {
			return mcp.NewToolResultError("Invalid UploadToken"), nil
		}
	}

	imagesDir, err := validateAndResolvePath(common.GetUploadDir(), "images")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid images directory: %v", err)), nil
	}
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create images directory: %v", err)), nil
	}

	filePath, err := validateAndResolvePath(imagesDir, filename)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid filename: %v", err)), nil
	}
	if _, err := os.Stat(filePath); err == nil {
		return mcp.NewToolResultError("File already exists"), nil
	}

	// 写入图片文件
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to write file: %v", err)), nil
	}

	// 生成缩略图
	thumbFilePath, err := validateAndResolvePath(imagesDir, "thumb_"+filename)
	if err != nil {
		os.Remove(filePath)
		return mcp.NewToolResultError(fmt.Sprintf("Invalid thumbnail filename: %v", err)), nil
	}
	thumbImage, err := imaging.Open(filePath)
	if err != nil {
		os.Remove(filePath) // 清理无效文件
		return mcp.NewToolResultError(fmt.Sprintf("Invalid image file, failed to parse: %v", err)), nil
	}
	thumbImage = imaging.Resize(thumbImage, 250, 150, imaging.Lanczos)
	if err := imaging.Save(thumbImage, thumbFilePath); err != nil {
		os.Remove(filePath)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to generate thumbnail: %v", err)), nil
	}

	// 获取文件信息
	fileInfo, _ := os.Stat(filePath)
	fileMD5, _ := utils.CalculateMD5(filePath)
	fileUUID := uuid.NewV4().String()

	// 写入 SQLite 数据库
	if common.EnableSqlite {
		go func() {
			db, err := common.GetDB()
			if err != nil {
				ylog.Errorf("MCP", "Failed to get DB for image record: %v", err)
				return
			}
			image := models.UploadImageModel{
				FileUUID:      fileUUID,
				Size:          fileInfo.Size(),
				FileName:      filename,
				FilePath:      filePath,
				ThumbFilePath: thumbFilePath,
				FileMD5:       fileMD5,
				DownloadCount: 0,
				Sort:          1000,
				UploadTime:    time.Now().Format("2006-01-02 15:04:05"),
				UploadIP:      "MCP",
				UploadUser:    "mcp",
				Status:        "done",
			}
			db.Create(&image)
		}()
	}

	ylog.Infof("MCP", "Image uploaded via MCP: %s (%s)", filename, utils.FormatSize(fileInfo.Size()))

	type UploadImageResult struct {
		FileUUID string `json:"file_uuid"`
		Filename string `json:"filename"`
		Size     string `json:"size"`
		MD5      string `json:"md5"`
		URL      string `json:"url"`
		ThumbURL string `json:"thumb_url"`
	}

	uploadResult := UploadImageResult{
		FileUUID: fileUUID,
		Filename: filename,
		Size:     utils.FormatSize(fileInfo.Size()),
		MD5:      fileMD5,
		URL:      "/api/v1/imageManage/download?filename=" + filename,
		ThumbURL: "/api/v1/imageManage/download?filename=thumb_" + filename,
	}

	resultData, _ := json.MarshalIndent(uploadResult, "", "  ")
	return mcp.NewToolResultText(string(resultData)), nil
}

// handleGetStatistics 处理统计信息查询
func handleGetStatistics(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	statsType := "all"
	if t, ok := args["type"].(string); ok && t != "" {
		statsType = t
	}

	type Statistics struct {
		Uploads   map[string]interface{} `json:"uploads,omitempty"`
		Downloads map[string]interface{} `json:"downloads,omitempty"`
		DiskUsage map[string]interface{} `json:"disk_usage,omitempty"`
	}

	stats := Statistics{}

	if statsType == "all" || statsType == "upload" {
		// 获取上传统计
		uploadDir := common.GetUploadDir()
		var totalSize int64
		var fileCount int

		filepath.Walk(uploadDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				totalSize += info.Size()
				fileCount++
			}
			return nil
		})

		stats.Uploads = map[string]interface{}{
			"file_count": fileCount,
			"total_size": utils.FormatSize(totalSize),
		}
	}

	if statsType == "all" || statsType == "disk" {
		usage, _ := getDiskUsage(common.GetUploadDir())
		stats.DiskUsage = usage
	}

	result, _ := json.MarshalIndent(stats, "", "  ")
	return mcp.NewToolResultText(string(result)), nil
}

// handleVerifyFileMD5 处理 MD5 验证
func handleVerifyFileMD5(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	filename, ok := args["filename"].(string)
	if !ok || filename == "" {
		return mcp.NewToolResultError("filename is required"), nil
	}

	expectedMD5, ok := args["expected_md5"].(string)
	if !ok || expectedMD5 == "" {
		return mcp.NewToolResultError("expected_md5 is required"), nil
	}

	// 使用安全的路径验证
	filePath, err := validateAndResolvePath(common.GetDownloadDir(), filename)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid filename: %v", err)), nil
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return mcp.NewToolResultError("File not found"), nil
	}

	actualMD5, err := utils.CalculateMD5(filePath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to calculate MD5: %v", err)), nil
	}

	type VerifyResult struct {
		Filename    string `json:"filename"`
		ExpectedMD5 string `json:"expected_md5"`
		ActualMD5   string `json:"actual_md5"`
		Match       bool   `json:"match"`
	}

	result := VerifyResult{
		Filename:    filename,
		ExpectedMD5: expectedMD5,
		ActualMD5:   actualMD5,
		Match:       strings.EqualFold(expectedMD5, actualMD5),
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

// ==================== Resource Handlers ====================

// handleFileListResource 文件列表资源处理器
func handleFileListResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	dirPath := common.GetDownloadDir()

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	type fileInfo struct {
		Name         string `json:"name"`
		Size         string `json:"size"`
		LastModified string `json:"last_modified"`
	}

	var fileList []fileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil || file.IsDir() {
			continue
		}
		fileList = append(fileList, fileInfo{
			Name:         file.Name(),
			Size:         utils.FormatSize(info.Size()),
			LastModified: info.ModTime().Format("2006-01-02 15:04:05"),
		})
	}

	data, _ := json.Marshal(fileList)
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "filelist://current",
			MIMEType: "application/json",
			Text:     string(data),
		},
	}, nil
}

// handleDiskUsageResource 磁盘使用资源处理器
func handleDiskUsageResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	usage, err := getDiskUsage(common.GetUploadDir())
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(usage)
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "disk://usage",
			MIMEType: "application/json",
			Text:     string(data),
		},
	}, nil
}

// handleSystemInfoResource 系统信息资源处理器
func handleSystemInfoResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	type SystemInfo struct {
		Version           string `json:"version"`
		UploadDir         string `json:"upload_dir"`
		DownloadDir       string `json:"download_dir"`
		P2PEnabled        bool   `json:"p2p_enabled"`
		SQLiteEnabled     bool   `json:"sqlite_enabled"`
		UploadTokenEnable bool   `json:"upload_token_enable"`
	}

	info := SystemInfo{
		Version:           common.Version,
		UploadDir:         common.GetUploadDir(),
		DownloadDir:       common.GetDownloadDir(),
		P2PEnabled:        common.P2pEnable,
		SQLiteEnabled:     common.EnableSqlite,
		UploadTokenEnable: common.EnableUploadToken,
	}

	data, _ := json.Marshal(info)
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "system://info",
			MIMEType: "application/json",
			Text:     string(data),
		},
	}, nil
}

// ==================== Helper Functions ====================

// getDiskUsage 获取磁盘使用情况
func getDiskUsage(path string) (map[string]interface{}, error) {
	// 这里简化处理，实际应该使用 gopsutil
	var totalSize, fileCount int64

	err := filepath.Walk(path, func(filepath string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			totalSize += info.Size()
			fileCount++
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"path":        path,
		"file_count":  fileCount,
		"total_size":  utils.FormatSize(totalSize),
		"total_bytes": totalSize,
	}, nil
}

// formatBytes 格式化字节数为人类可读格式
func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)
	switch {
	case bytes >= int64(TB):
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(TB))
	case bytes >= int64(GB):
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= int64(MB):
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= int64(KB):
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	}
	return fmt.Sprintf("%d B", bytes)
}

// ==================== v0.5.0 New Tool Handlers ====================

// handleCreateFolder 处理创建文件夹
func handleCreateFolder(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	name, ok := args["name"].(string)
	if !ok || name == "" {
		return mcp.NewToolResultError("name is required"), nil
	}

	// 安全校验文件夹名
	safeName, err := common.NormalizeSafeFileName(name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid folder name: %v", err)), nil
	}
	if strings.HasPrefix(safeName, ".") {
		return mcp.NewToolResultError("Invalid folder name: hidden folders not allowed"), nil
	}

	// 解析父目录
	basePath := common.GetDownloadDir()
	if dirVal, ok := args["dir"].(string); ok && dirVal != "" {
		basePath, err = validateAndResolvePath(common.GetDownloadDir(), dirVal)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid directory path: %v", err)), nil
		}
	}

	folderPath, err := validateAndResolvePath(basePath, safeName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid folder path: %v", err)), nil
	}

	// 检查是否已存在
	if _, err := os.Stat(folderPath); err == nil {
		return mcp.NewToolResultError("Folder already exists"), nil
	}

	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create folder: %v", err)), nil
	}

	ylog.Infof("MCP", "Folder created via MCP: %s", folderPath)

	type CreateFolderResult struct {
		Name    string `json:"name"`
		Path    string `json:"path"`
		Message string `json:"message"`
	}

	result := CreateFolderResult{
		Name:    safeName,
		Path:    folderPath,
		Message: "Folder created successfully",
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

// handleRenameFile 处理重命名文件或文件夹
func handleRenameFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	oldName, ok := args["old_name"].(string)
	if !ok || oldName == "" {
		return mcp.NewToolResultError("old_name is required"), nil
	}

	newName, ok := args["new_name"].(string)
	if !ok || newName == "" {
		return mcp.NewToolResultError("new_name is required"), nil
	}

	// 安全校验
	safeOldName, err := common.NormalizeSafeFileName(oldName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid old name: %v", err)), nil
	}
	safeNewName, err := common.NormalizeSafeFileName(newName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid new name: %v", err)), nil
	}
	if strings.HasPrefix(safeNewName, ".") {
		return mcp.NewToolResultError("Invalid new name: hidden files not allowed"), nil
	}

	// 解析目录
	basePath := common.GetDownloadDir()
	if dirVal, ok := args["dir"].(string); ok && dirVal != "" {
		basePath, err = validateAndResolvePath(common.GetDownloadDir(), dirVal)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid directory path: %v", err)), nil
		}
	}

	oldPath, err := validateAndResolvePath(basePath, safeOldName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid old path: %v", err)), nil
	}
	newPath, err := validateAndResolvePath(basePath, safeNewName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid new path: %v", err)), nil
	}

	// 检查源文件是否存在
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return mcp.NewToolResultError("Source file not found"), nil
	}

	// 检查目标是否已存在
	if _, err := os.Stat(newPath); err == nil {
		return mcp.NewToolResultError("Target name already exists"), nil
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to rename: %v", err)), nil
	}

	ylog.Infof("MCP", "Renamed via MCP: %s -> %s", oldPath, newPath)

	type RenameResult struct {
		OldName string `json:"old_name"`
		NewName string `json:"new_name"`
		Message string `json:"message"`
	}

	result := RenameResult{
		OldName: safeOldName,
		NewName: safeNewName,
		Message: "Renamed successfully",
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

// handleBatchDeleteFiles 处理批量删除文件
func handleBatchDeleteFiles(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	filesStr, ok := args["files"].(string)
	if !ok || filesStr == "" {
		return mcp.NewToolResultError("files is required (JSON array of filenames)"), nil
	}

	// 解析 JSON 数组
	var files []string
	if err := json.Unmarshal([]byte(filesStr), &files); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid files format, expected JSON array: %v", err)), nil
	}

	if len(files) == 0 {
		return mcp.NewToolResultError("files list is empty"), nil
	}

	// 解析目录
	basePath := common.GetDownloadDir()
	if dirVal, ok := args["dir"].(string); ok && dirVal != "" {
		var err error
		basePath, err = validateAndResolvePath(common.GetDownloadDir(), dirVal)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid directory path: %v", err)), nil
		}
	}

	type deleteResult struct {
		Deleted []string            `json:"deleted"`
		Failed  []map[string]string `json:"failed"`
	}

	result := deleteResult{}

	for _, fileName := range files {
		safeName, err := common.NormalizeSafeFileName(fileName)
		if err != nil {
			result.Failed = append(result.Failed, map[string]string{"file": fileName, "error": "invalid filename"})
			continue
		}

		filePath, err := validateAndResolvePath(basePath, safeName)
		if err != nil {
			result.Failed = append(result.Failed, map[string]string{"file": fileName, "error": "invalid path"})
			continue
		}

		info, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				result.Failed = append(result.Failed, map[string]string{"file": fileName, "error": "file not found"})
			} else {
				result.Failed = append(result.Failed, map[string]string{"file": fileName, "error": err.Error()})
			}
			continue
		}

		if info.IsDir() {
			// 删除目录需要目录为空
			if err := os.Remove(filePath); err != nil {
				result.Failed = append(result.Failed, map[string]string{"file": fileName, "error": "directory is not empty or cannot be removed"})
				continue
			}
		} else {
			if err := os.Remove(filePath); err != nil {
				result.Failed = append(result.Failed, map[string]string{"file": fileName, "error": err.Error()})
				continue
			}
		}

		result.Deleted = append(result.Deleted, fileName)
		ylog.Infof("MCP", "Batch delete via MCP: %s", filePath)
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

// handleGetFileOverview 处理文件总览统计
func handleGetFileOverview(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	baseDir := common.GetDownloadDir()

	var totalFiles int64
	var totalDirs int64
	var totalSize int64

	// 递归统计
	filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 跳过无法访问的文件
		}
		if info.IsDir() {
			totalDirs++
		} else {
			totalFiles++
			totalSize += info.Size()
		}
		return nil
	})

	// 减去根目录本身
	if totalDirs > 0 {
		totalDirs--
	}

	type FileOverview struct {
		TotalFiles         int64  `json:"total_files"`
		TotalDirs          int64  `json:"total_dirs"`
		TotalSize          int64  `json:"total_size"`
		TotalSizeFormatted string `json:"total_size_formatted"`
	}

	overview := FileOverview{
		TotalFiles:         totalFiles,
		TotalDirs:          totalDirs,
		TotalSize:          totalSize,
		TotalSizeFormatted: formatBytes(totalSize),
	}

	data, _ := json.MarshalIndent(overview, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

// handleGetDownloadHistory 处理下载历史日志查询
func handleGetDownloadHistory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !common.EnableSqlite {
		return mcp.NewToolResultError("SQLite is not enabled"), nil
	}

	args := request.GetArguments()

	page := 1
	pageSize := 20
	filename := ""
	fileMD5 := ""
	ip := ""

	if v, ok := args["page"].(float64); ok {
		page = int(v)
	}
	if v, ok := args["page_size"].(float64); ok {
		pageSize = int(v)
	}
	if v, ok := args["filename"].(string); ok {
		filename = v
	}
	if v, ok := args["file_md5"].(string); ok {
		fileMD5 = v
	}
	if v, ok := args["ip"].(string); ok {
		ip = v
	}

	db, err := common.GetDB()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get database: %v", err)), nil
	}

	offset := (page - 1) * pageSize
	query := db.Model(&common.DownloadLogModel{}).Offset(offset).Limit(pageSize).Order("download_time DESC")
	countQuery := db.Model(&common.DownloadLogModel{})

	if filename != "" {
		query = query.Where("filename LIKE ?", "%"+filename+"%")
		countQuery = countQuery.Where("filename LIKE ?", "%"+filename+"%")
	}
	if fileMD5 != "" {
		query = query.Where("file_md5 = ?", fileMD5)
		countQuery = countQuery.Where("file_md5 = ?", fileMD5)
	}
	if ip != "" {
		query = query.Where("ip = ?", ip)
		countQuery = countQuery.Where("ip = ?", ip)
	}

	var logs []common.DownloadLogModel
	if err := query.Find(&logs).Error; err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to query download history: %v", err)), nil
	}

	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to count download history: %v", err)), nil
	}

	// 转换为 map 返回
	var logList []map[string]interface{}
	for _, log := range logs {
		logList = append(logList, map[string]interface{}{
			"id":            log.ID,
			"ip":            log.IP,
			"appkey":        log.AppKey,
			"download_time": log.DownloadTime,
			"filename":      log.FileName,
			"file_size":     log.FileSize,
			"file_md5":      log.FileMD5,
		})
	}

	response := map[string]interface{}{
		"list":      logList,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	}

	data, _ := json.MarshalIndent(response, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}
