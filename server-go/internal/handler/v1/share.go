package v1

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"httpcat/internal/common"
	"httpcat/internal/common/ylog"
	"httpcat/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

// ===== 请求/响应结构 =====

type CreateShareRequest struct {
	FilePath     string `json:"filePath" binding:"required"`     // 文件相对路径（如 "test.txt"）
	FileName     string `json:"fileName" binding:"required"`     // 文件名
	FileType     string `json:"fileType"`                        // file / image
	ExtractCode  string `json:"extractCode"`                     // 提取码（空=不设置）
	ExpireHours  int    `json:"expireHours"`                     // 过期时间（小时），0=永不过期
	MaxDownloads int    `json:"maxDownloads"`                    // 最大下载次数，0=不限
}

type ShareInfoResponse struct {
	ShareCode    string     `json:"shareCode"`
	FileName     string     `json:"fileName"`
	FileType     string     `json:"fileType"`
	HasExtract   bool       `json:"hasExtractCode"`  // 是否需要提取码
	ExpireAt     *time.Time `json:"expireAt"`
	MaxDownloads int        `json:"maxDownloads"`
	CurDownloads int        `json:"curDownloads"`
	IsActive     bool       `json:"isActive"`
	CreatedBy    string     `json:"createdBy"`
	CreatedAt    time.Time  `json:"createdAt"`
}

// ===== 工具函数 =====

const shareCodeChars = "abcdefghijkmnpqrstuvwxyz23456789"

func generateShareCode(length int) string {
	code := make([]byte, length)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(shareCodeChars))))
		code[i] = shareCodeChars[n.Int64()]
	}
	return string(code)
}

func generateExtractCode() string {
	code := make([]byte, 4)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(shareCodeChars))))
		code[i] = shareCodeChars[n.Int64()]
	}
	return string(code)
}

func isShareValid(share *models.ShareModel) (bool, string) {
	if !share.IsActive {
		return false, "该分享已被取消"
	}
	if share.ExpireAt != nil && time.Now().After(*share.ExpireAt) {
		return false, "该分享链接已过期"
	}
	if share.MaxDownloads > 0 && share.CurDownloads >= share.MaxDownloads {
		return false, "该分享已达到最大下载次数"
	}
	return true, ""
}

// ===== Handler 实现 =====

// checkShareAnonymousAccess 检查分享公开路由的匿名访问权限
// 如果 ShareAnonymousAccess 为 false，则需要登录（有有效 JWT Token）
func checkShareAnonymousAccess(c *gin.Context) bool {
	if common.ShareAnonymousAccess {
		return true // 允许匿名访问
	}

	// 不允许匿名访问，检查 JWT Token
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "此分享链接需要登录后才能访问"})
		return false
	}
	if len(token) > 7 && strings.ToUpper(token[0:7]) == "BEARER " {
		token = token[7:]
	}

	jwtSecret := []byte(common.JwtSecret)
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil || !jwtToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "此分享链接需要登录后才能访问"})
		return false
	}

	return true
}

// CreateShare 创建分享 POST /api/v1/share
func CreateShare(c *gin.Context) {
	var req CreateShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	username, _ := c.Get("user")

	// 验证文件存在
	fileType := req.FileType
	if fileType == "" {
		fileType = "file"
	}

	var baseDir string
	if fileType == "image" {
		imgDir, err := common.ResolvePathWithinBase(common.GetUploadDir(), "images")
		if err != nil {
			common.BadRequest(c, "无效的图片目录")
			return
		}
		baseDir = imgDir
	} else {
		baseDir = common.GetDownloadDir()
	}

	filePath, err := common.ResolvePathWithinBase(baseDir, req.FilePath)
	if err != nil {
		common.BadRequest(c, "无效的文件路径")
		return
	}

	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		common.BadRequest(c, "文件不存在")
		return
	}

	db, err := common.GetDB()
	if err != nil {
		common.CreateResponse(c, common.ErrorCode, "数据库连接失败")
		return
	}

	// 生成唯一分享码（重试3次防碰撞）
	var shareCode string
	for i := 0; i < 3; i++ {
		shareCode = generateShareCode(8)
		var count int64
		db.Model(&models.ShareModel{}).Where("share_code = ?", shareCode).Count(&count)
		if count == 0 {
			break
		}
	}

	// 提取码
	extractCode := strings.TrimSpace(req.ExtractCode)
	if extractCode == "auto" {
		extractCode = generateExtractCode()
	}

	// 过期时间
	var expireAt *time.Time
	if req.ExpireHours > 0 {
		t := time.Now().Add(time.Duration(req.ExpireHours) * time.Hour)
		expireAt = &t
	}

	share := models.ShareModel{
		ShareCode:    shareCode,
		FilePath:     req.FilePath,
		FileName:     req.FileName,
		FileType:     fileType,
		CreatedBy:    username.(string),
		ExtractCode:  extractCode,
		ExpireAt:     expireAt,
		MaxDownloads: req.MaxDownloads,
		CurDownloads: 0,
		IsActive:     true,
	}

	if err := db.Create(&share).Error; err != nil {
		ylog.Errorf("CreateShare", "创建分享失败: %v", err)
		common.CreateResponse(c, common.ErrorCode, "创建分享失败")
		return
	}

	common.CreateResponse(c, common.SuccessCode, gin.H{
		"shareCode":   shareCode,
		"extractCode": extractCode,
		"shareUrl":    fmt.Sprintf("/s/%s", shareCode),
		"expireAt":    expireAt,
	})
}

// ListShares 我的分享列表 GET /api/v1/shares
func ListShares(c *gin.Context) {
	username, _ := c.Get("user")

	db, err := common.GetDB()
	if err != nil {
		common.CreateResponse(c, common.ErrorCode, "数据库连接失败")
		return
	}

	pageStr := c.DefaultQuery("current", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var shares []models.ShareModel
	var total int64

	query := db.Model(&models.ShareModel{}).Where("created_by = ?", username.(string))
	query.Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&shares)

	common.CreateResponse(c, common.SuccessCode, gin.H{
		"list":     shares,
		"current":  page,
		"pageSize": pageSize,
		"total":    total,
	})
}

// DeleteShare 取消分享 DELETE /api/v1/share/:code
func DeleteShare(c *gin.Context) {
	code := c.Param("code")
	username, _ := c.Get("user")

	db, err := common.GetDB()
	if err != nil {
		common.CreateResponse(c, common.ErrorCode, "数据库连接失败")
		return
	}

	result := db.Model(&models.ShareModel{}).
		Where("share_code = ? AND created_by = ?", code, username.(string)).
		Update("is_active", false)

	if result.RowsAffected == 0 {
		common.BadRequest(c, "分享不存在或无权操作")
		return
	}

	common.CreateResponse(c, common.SuccessCode, "分享已取消")
}

// GetShareInfo 分享落地页信息 GET /s/:code
func GetShareInfo(c *gin.Context) {
	if !checkShareAnonymousAccess(c) {
		return
	}

	code := c.Param("code")

	db, err := common.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务暂不可用"})
		return
	}

	var share models.ShareModel
	if err := db.Where("share_code = ?", code).First(&share).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分享不存在"})
		return
	}

	valid, reason := isShareValid(&share)
	if !valid {
		c.JSON(http.StatusOK, gin.H{
			"valid":  false,
			"reason": reason,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"share": ShareInfoResponse{
			ShareCode:    share.ShareCode,
			FileName:     share.FileName,
			FileType:     share.FileType,
			HasExtract:   share.ExtractCode != "",
			ExpireAt:     share.ExpireAt,
			MaxDownloads: share.MaxDownloads,
			CurDownloads: share.CurDownloads,
			IsActive:     share.IsActive,
			CreatedBy:    share.CreatedBy,
			CreatedAt:    share.CreatedAt,
		},
	})
}

// VerifyShareCode 验证提取码 POST /s/:code/verify
func VerifyShareCode(c *gin.Context) {
	if !checkShareAnonymousAccess(c) {
		return
	}

	code := c.Param("code")

	var req struct {
		ExtractCode string `json:"extractCode" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, "请输入提取码")
		return
	}

	db, err := common.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务暂不可用"})
		return
	}

	var share models.ShareModel
	if err := db.Where("share_code = ?", code).First(&share).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分享不存在"})
		return
	}

	valid, reason := isShareValid(&share)
	if !valid {
		c.JSON(http.StatusOK, gin.H{"valid": false, "reason": reason})
		return
	}

	if share.ExtractCode != "" && share.ExtractCode != strings.TrimSpace(req.ExtractCode) {
		c.JSON(http.StatusOK, gin.H{"valid": false, "reason": "提取码错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true})
}

// DownloadShareFile 下载分享文件 GET /s/:code/download
func DownloadShareFile(c *gin.Context) {
	if !checkShareAnonymousAccess(c) {
		return
	}

	code := c.Param("code")
	extractCode := c.Query("code") // 提取码通过 query 传递

	db, err := common.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务暂不可用"})
		return
	}

	var share models.ShareModel
	if err := db.Where("share_code = ?", code).First(&share).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分享不存在"})
		return
	}

	valid, reason := isShareValid(&share)
	if !valid {
		c.JSON(http.StatusForbidden, gin.H{"error": reason})
		return
	}

	// 验证提取码
	if share.ExtractCode != "" && share.ExtractCode != strings.TrimSpace(extractCode) {
		c.JSON(http.StatusForbidden, gin.H{"error": "提取码错误"})
		return
	}

	// 解析文件路径
	var baseDir string
	if share.FileType == "image" {
		imgDir, err := common.ResolvePathWithinBase(common.GetUploadDir(), "images")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文件目录无效"})
			return
		}
		baseDir = imgDir
	} else {
		baseDir = common.GetDownloadDir()
	}

	filePath, err := common.ResolvePathWithinBase(baseDir, share.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件路径无效"})
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		ylog.Errorf("DownloadShareFile", "打开文件失败: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil || fileInfo.IsDir() {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 原子更新下载计数（防并发超发）
	// 使用 SQL 条件确保在有限制时不会超过 max_downloads
	updateResult := db.Model(&models.ShareModel{}).
		Where("share_code = ? AND is_active = ? AND (max_downloads = 0 OR cur_downloads < max_downloads)",
			share.ShareCode, true).
		Update("cur_downloads", gorm.Expr("cur_downloads + 1"))

	if updateResult.RowsAffected == 0 {
		// 可能在请求间隙被其他请求抢完了
		c.JSON(http.StatusForbidden, gin.H{"error": "该分享已达到最大下载次数"})
		return
	}

	// 发送文件
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(share.FileName))
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	buf := make([]byte, 4096)
	for {
		n, err := file.Read(buf)
		if n == 0 {
			break
		}
		if err != nil && err != io.EOF {
			ylog.Errorf("DownloadShareFile", "读取文件失败: %v", err)
			break
		}
		_, _ = c.Writer.Write(buf[:n])
	}
}

// ShareStats 分享统计 GET /api/v1/share/stats
func ShareStats(c *gin.Context) {
	username, _ := c.Get("user")

	db, err := common.GetDB()
	if err != nil {
		common.CreateResponse(c, common.ErrorCode, "数据库连接失败")
		return
	}

	var totalShares int64
	var activeShares int64
	var totalDownloads int64
	var expiredShares int64

	now := time.Now()

	db.Model(&models.ShareModel{}).Where("created_by = ?", username.(string)).Count(&totalShares)
	db.Model(&models.ShareModel{}).Where("created_by = ? AND is_active = ?", username.(string), true).Count(&activeShares)

	// 过期的分享数
	db.Model(&models.ShareModel{}).
		Where("created_by = ? AND ((expire_at IS NOT NULL AND expire_at < ?) OR (is_active = ? AND max_downloads > 0 AND cur_downloads >= max_downloads))",
			username.(string), now, false).
		Count(&expiredShares)

	// 总下载量
	var result struct {
		Total int64
	}
	db.Model(&models.ShareModel{}).
		Select("COALESCE(SUM(cur_downloads), 0) as total").
		Where("created_by = ?", username.(string)).
		Scan(&result)
	totalDownloads = result.Total

	common.CreateResponse(c, common.SuccessCode, gin.H{
		"totalShares":    totalShares,
		"activeShares":   activeShares,
		"expiredShares":  expiredShares,
		"totalDownloads": totalDownloads,
	})
}

// GetShareConfig 获取分享功能配置 GET /api/v1/share/config
func GetShareConfig(c *gin.Context) {
	common.CreateResponse(c, common.SuccessCode, gin.H{
		"shareEnable":     common.ShareEnable,
		"anonymousAccess": common.ShareAnonymousAccess,
	})
}
