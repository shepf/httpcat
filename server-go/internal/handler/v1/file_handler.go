package v1

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"httpcat/internal/common"
	"httpcat/internal/common/utils"
	"httpcat/internal/common/ylog"
	"httpcat/internal/models"
	"httpcat/internal/storage"
	"httpcat/internal/storage/auth"

	"github.com/gin-gonic/gin"
)

// GetDirConf 获取配置文件中的上传下载目录配置
func GetDirConf(c *gin.Context) {
	dirConf := map[string]string{
		"UploadDir":  common.GetUploadDir(),
		"DownloadDir": common.GetDownloadDir(),
		"StaticDir":  common.StaticDir,
	}

	common.CreateResponse(c, common.SuccessCode, dirConf)
}

// UploadFile 处理文件上传
func UploadFile(c *gin.Context) {
	if !common.FileUploadEnable {
		common.CreateResponse(c, common.ErrorCode, "File service is not enabled")
		return
	}

	file, header, err := c.Request.FormFile("f1")
	if err != nil {
		common.BadRequest(c, "Bad request,check your file~")
		return
	}
	defer file.Close()

	appkey := ""
	if common.EnableUploadToken {
		uploadToken := c.Request.Header.Get("UploadToken")
		if uploadToken == "" {
			common.BadRequest(c, "UploadToken is empty")
			return
		}

		parts := strings.Split(uploadToken, ":")
		if len(parts) != 3 {
			common.Unauthorized(c, "Invalid UploadToken format")
			return
		}
		appkey = parts[0]
		common.UploadTokenLock.RLock()
		tokenItem, ok := common.UploadTokenTable[appkey]
		common.UploadTokenLock.RUnlock()
		if !ok {
			common.Unauthorized(c, "Invalid Appkey")
			return
		}

		if tokenItem.State == "closed" {
			common.Unauthorized(c, "Invalid Appkey, appkey is closed")
			return
		}

		mac := auth.New(appkey, tokenItem.Appsecret)
		if !mac.VerifyUploadToken(uploadToken) {
			common.Unauthorized(c, "UploadToken is invalid")
			return
		}
	}

	filename, err := common.NormalizeSafeFileName(header.Filename)
	if err != nil {
		common.BadRequest(c, "invalid filename")
		return
	}

	if err := os.MkdirAll(common.GetUploadDir(), 0755); err != nil {
		ylog.Errorf("uploadFile", "创建目录失败", err)
		common.CreateResponse(c, common.ErrorCode, "Failed to prepare upload directory")
		return
	}

	filePath, err := common.ResolvePathWithinBase(common.GetUploadDir(), filename)
	if err != nil {
		common.BadRequest(c, "invalid filename")
		return
	}

	ylog.Infof("uploadFile", "upload file to: %s", filePath)
	out, err := os.Create(filePath)
	if err != nil {
		ylog.Errorf("uploadFile", "创建文件失败: %v", err)
		common.CreateResponse(c, common.ErrorCode, "Failed to create file")
		return
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		ylog.Errorf("uploadFile", "写入文件失败: %v", err)
		common.CreateResponse(c, common.ErrorCode, "Failed to save file")
		return
	}

	ip := c.ClientIP()
	uploadTime := time.Now().Format("2006-01-02 15:04:05")
	fileInfo, _ := os.Stat(filePath)
	fileSize := utils.FormatSize(fileInfo.Size())
	fileMD5, _ := utils.CalculateMD5(filePath)
	fileCreatedTime := fileInfo.ModTime().Unix()
	fileModifiedTime := fileInfo.ModTime().Unix()

	fmt.Println("PersistentNotifyURL:", common.PersistentNotifyURL)
	if common.PersistentNotifyURL != "" {
		ylog.Infof("uploadFile", "send notify to: %s", common.PersistentNotifyURL)

		markdownContent := fmt.Sprintf(`>有文件上传归档,上传信息：
			- IP地址：%s
			- 上传时间：%s
			- 文件名：%s
			- 文件大小：%s
			- 文件MD5：%s`, ip, uploadTime, filename, fileSize, fileMD5)
		ylog.Infof("uploadFile", "markdownContent:%s", markdownContent)

		go utils.SendNotify(common.PersistentNotifyURL, markdownContent)
	}

	if common.EnableSqlite {
		ylog.Infof("uploadFile", "sqliteInsert enable")
		go insertUploadLog(ip, appkey, uploadTime, filename, fileSize, fileMD5, fileCreatedTime, fileModifiedTime)
	}

	common.CreateResponse(c, common.SuccessCode, "upload successful!")
}

// insertUploadLog 使用 GORM 单例插入上传日志（替代旧的 sqliteInsert 使用 raw sql.Open）
func insertUploadLog(ip string, appkey string, uploadTime string, filename string, fileSize string, fileMD5 string,
	fileCreatedTime int64, fileModifiedTime int64) {
	ylog.Infof("insertUploadLog", "start")

	db, err := common.GetDB()
	if err != nil {
		ylog.Errorf("insertUploadLog", "获取数据库连接失败: %v", err)
		return
	}

	log := models.UploadLogModel{
		IP:               ip,
		Appkey:           appkey,
		UploadTime:       uploadTime,
		FileName:         filename,
		FileSize:         fileSize,
		FileMD5:          fileMD5,
		FileCreatedTime:  fileCreatedTime,
		FileModifiedTime: fileModifiedTime,
	}

	if err := db.Create(&log).Error; err != nil {
		ylog.Errorf("insertUploadLog", "插入上传日志失败: %v", err)
		return
	}

	ylog.Infof("insertUploadLog", "end")
}

// ListFiles 获取目录文件列表（支持子目录导航，包含目录条目）
func ListFiles(c *gin.Context) {
	dirPath, err := common.ResolvePathWithinBase(common.GetDownloadDir(), c.Query("dir"))
	if err != nil {
		common.BadRequest(c, "invalid dir")
		return
	}
	ylog.Infof("ListFiles", "dirPath:%s", dirPath)

	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			ylog.Errorf("ListFiles", "目录不存在: %v", err)
			common.CreateResponse(c, common.DirISNotExists, "Directory does not exist")
		} else {
			ylog.Errorf("ListFiles", "读取目录失败: %v", err)
			common.CreateResponse(c, common.ReadDirFailed, "Failed to read the directory")
		}
		c.AbortWithStatus(500)
		return
	}
	if !info.IsDir() {
		common.BadRequest(c, "dir must be a directory")
		return
	}

	// 使用 os.ReadDir 替代废弃的 ioutil.ReadDir
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			ylog.Errorf("ListFiles", "目录不存在: %v", err)
			common.CreateResponse(c, common.DirISNotExists, "Directory does not exist")
		} else {
			ylog.Errorf("ListFiles", "读取目录失败: %v", err)
			common.CreateResponse(c, common.ReadDirFailed, "Failed to read the directory")
		}
		c.AbortWithStatus(500)
		return
	}

	// 获取 FileInfo 列表用于排序
	type fileItem struct {
		entry os.DirEntry
		info  os.FileInfo
	}
	var items []fileItem
	for _, entry := range entries {
		fi, err := entry.Info()
		if err != nil {
			continue
		}
		items = append(items, fileItem{entry: entry, info: fi})
	}

	// 目录优先，然后按修改时间倒序
	sort.SliceStable(items, func(i, j int) bool {
		iIsDir := items[i].info.IsDir()
		jIsDir := items[j].info.IsDir()
		if iIsDir != jIsDir {
			return iIsDir // 目录排在前面
		}
		return items[j].info.ModTime().Before(items[i].info.ModTime())
	})

	var fileList []map[string]interface{}
	for _, item := range items {
		fileEntry := map[string]interface{}{
			"FileName":     item.info.Name(),
			"LastModified": item.info.ModTime().Format("2006-01-02 15:04:05"),
			"IsDir":        item.info.IsDir(),
		}
		if item.info.IsDir() {
			fileEntry["Size"] = "-"
		} else {
			fileEntry["Size"] = utils.FormatSize(item.info.Size())
		}
		fileList = append(fileList, fileEntry)
	}

	common.CreateResponse(c, common.SuccessCode, fileList)
}

// GetFileInfo 获取某个文件的详细信息
func GetFileInfo(c *gin.Context) {
	fileName := c.Query("filename")
	fileMD5Param := c.Query("file_md5")

	if fileName == "" {
		common.BadRequest(c, "filename is required")
		return
	}

	filePath, err := common.ResolvePathWithinBase(common.GetDownloadDir(), fileName)
	if err != nil {
		common.BadRequest(c, "invalid filename")
		return
	}
	ylog.Infof("GetFileInfo", "filePath:%s", filePath)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			ylog.Errorf("GetFileInfo", "文件不存在: %v", err)
			common.CreateResponse(c, common.FileIsNotExists, nil)
		} else {
			ylog.Errorf("GetFileInfo", "获取文件信息失败: %v", err)
			common.CreateResponse(c, common.ErrorCode, "Failed to get file information")
		}
		c.AbortWithStatus(500)
		return
	}
	if fileInfo.IsDir() {
		common.BadRequest(c, "invalid filename")
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		ylog.Errorf("GetFileInfo", "打开文件失败: %v", err)
		common.CreateResponse(c, common.ErrorCode, "Failed to open the file")
		c.AbortWithStatus(500)
		return
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		c.AbortWithStatus(500)
		return
	}

	md5Hash := hex.EncodeToString(hash.Sum(nil))

	fileEntry := map[string]interface{}{
		"fileName":     fileInfo.Name(),
		"lastModified": fileInfo.ModTime().Format("2006-01-02 15:04:05"),
		"size":         utils.FormatSize(fileInfo.Size()),
		"md5":          md5Hash,
		"md5Match":     fileMD5Param == "" || fileMD5Param == md5Hash,
	}

	common.CreateResponse(c, common.SuccessCode, fileEntry)
}

// UploadHistoryLogs 获取上传文件历史记录（分页）
func UploadHistoryLogs(c *gin.Context) {
	var params struct {
		Current  int    `form:"current" binding:"required"`
		PageSize int    `form:"pageSize" binding:"required"`
		FileName string `form:"filename"`
		FileMD5  string `form:"file_md5"`
		IP       string `form:"ip"`
		AppKey   string `form:"appkey"`
	}
	if err := c.ShouldBindQuery(&params); err != nil {
		ylog.Errorf("UploadHistoryLogs", "请求参数错误: %s", err.Error())
		common.CreateResponse(c, common.ParamInvalidErrorCode, err.Error())
		return
	}

	db, err := common.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	offset := (params.Current - 1) * params.PageSize
	query := db.Table("t_upload_log").Offset(offset).Limit(params.PageSize).Order("upload_time DESC")
	if params.FileName != "" {
		query = query.Where("filename LIKE ?", "%"+params.FileName+"%")
	}
	if params.FileMD5 != "" {
		query = query.Where("file_md5 = ?", params.FileMD5)
	}
	if params.IP != "" {
		query = query.Where("ip = ?", params.IP)
	}
	if params.AppKey != "" {
		query = query.Where("appkey = ?", params.AppKey)
	}

	var logs []models.UploadLogModel
	if err := query.Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 查询总数
	var total int64
	countQuery := db.Table("t_upload_log")
	if params.FileName != "" {
		countQuery = countQuery.Where("filename LIKE ?", "%"+params.FileName+"%")
	}
	if params.FileMD5 != "" {
		countQuery = countQuery.Where("file_md5 = ?", params.FileMD5)
	}
	if params.IP != "" {
		countQuery = countQuery.Where("ip = ?", params.IP)
	}
	if params.AppKey != "" {
		countQuery = countQuery.Where("appkey = ?", params.AppKey)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	common.CreateResponse(c, common.SuccessCode, gin.H{
		"list":     logs,
		"current":  params.Current,
		"pageSize": params.PageSize,
		"total":    total,
	})
}

// DeleteHistoryLogs 删除上传历史记录
func DeleteHistoryLogs(c *gin.Context) {
	ids, exists := c.GetQueryArray("id")
	if !exists || len(ids) == 0 {
		ylog.Errorf("DeleteHistoryLogs", "请求参数错误")
		common.CreateResponse(c, common.ParamInvalidErrorCode, "Invalid ID")
		return
	}
	ylog.Infof("DeleteHistoryLogs", "删除上传文件日志 ids: %v", ids)

	db, err := common.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.Table("t_upload_log").Where("id IN ?", ids).Delete(&models.UploadLogModel{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	common.CreateResponse(c, common.SuccessCode, nil)
}

// DeleteFiles 批量删除文件
func DeleteFiles(c *gin.Context) {
	var req struct {
		Files []string `json:"files" binding:"required"`
		Dir   string   `json:"dir"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, "invalid request body")
		return
	}

	if len(req.Files) == 0 {
		common.BadRequest(c, "files list is empty")
		return
	}

	basePath := common.GetDownloadDir()
	if req.Dir != "" {
		var err error
		basePath, err = common.ResolvePathWithinBase(common.GetDownloadDir(), req.Dir)
		if err != nil {
			common.BadRequest(c, "invalid dir")
			return
		}
	}

	var deleted []string
	var failed []map[string]string

	for _, fileName := range req.Files {
		safeName, err := common.NormalizeSafeFileName(fileName)
		if err != nil {
			failed = append(failed, map[string]string{"file": fileName, "error": "invalid filename"})
			continue
		}

		filePath, err := common.ResolvePathWithinBase(basePath, safeName)
		if err != nil {
			failed = append(failed, map[string]string{"file": fileName, "error": "invalid path"})
			continue
		}

		info, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				failed = append(failed, map[string]string{"file": fileName, "error": "file not found"})
			} else {
				failed = append(failed, map[string]string{"file": fileName, "error": err.Error()})
			}
			continue
		}

		if info.IsDir() {
			// 删除目录需要目录为空
			if err := os.Remove(filePath); err != nil {
				failed = append(failed, map[string]string{"file": fileName, "error": "directory is not empty or cannot be removed"})
				continue
			}
		} else {
			if err := os.Remove(filePath); err != nil {
				failed = append(failed, map[string]string{"file": fileName, "error": err.Error()})
				continue
			}
		}

		deleted = append(deleted, fileName)
		ylog.Infof("DeleteFiles", "deleted: %s", filePath)
	}

	common.CreateResponse(c, common.SuccessCode, gin.H{
		"deleted": deleted,
		"failed":  failed,
	})
}

// CreateFolder 创建文件夹
func CreateFolder(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Dir  string `json:"dir"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, "invalid request body")
		return
	}

	safeName, err := common.NormalizeSafeFileName(req.Name)
	if err != nil {
		common.BadRequest(c, "invalid folder name")
		return
	}

	basePath := common.GetDownloadDir()
	if req.Dir != "" {
		basePath, err = common.ResolvePathWithinBase(common.GetDownloadDir(), req.Dir)
		if err != nil {
			common.BadRequest(c, "invalid dir")
			return
		}
	}

	folderPath, err := common.ResolvePathWithinBase(basePath, safeName)
	if err != nil {
		common.BadRequest(c, "invalid folder path")
		return
	}

	// 检查是否已存在
	if _, err := os.Stat(folderPath); err == nil {
		common.CreateResponse(c, common.ErrorCode, "folder already exists")
		return
	}

	if err := os.MkdirAll(folderPath, 0755); err != nil {
		ylog.Errorf("CreateFolder", "创建文件夹失败: %v", err)
		common.CreateResponse(c, common.ErrorCode, "Failed to create folder")
		return
	}

	ylog.Infof("CreateFolder", "created folder: %s", folderPath)
	common.CreateResponse(c, common.SuccessCode, "folder created")
}

// RenameFile 重命名文件或文件夹
func RenameFile(c *gin.Context) {
	var req struct {
		OldName string `json:"oldName" binding:"required"`
		NewName string `json:"newName" binding:"required"`
		Dir     string `json:"dir"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, "invalid request body")
		return
	}

	safeOldName, err := common.NormalizeSafeFileName(req.OldName)
	if err != nil {
		common.BadRequest(c, "invalid old name")
		return
	}
	safeNewName, err := common.NormalizeSafeFileName(req.NewName)
	if err != nil {
		common.BadRequest(c, "invalid new name")
		return
	}

	basePath := common.GetDownloadDir()
	if req.Dir != "" {
		basePath, err = common.ResolvePathWithinBase(common.GetDownloadDir(), req.Dir)
		if err != nil {
			common.BadRequest(c, "invalid dir")
			return
		}
	}

	oldPath, err := common.ResolvePathWithinBase(basePath, safeOldName)
	if err != nil {
		common.BadRequest(c, "invalid old path")
		return
	}
	newPath, err := common.ResolvePathWithinBase(basePath, safeNewName)
	if err != nil {
		common.BadRequest(c, "invalid new path")
		return
	}

	// 检查源文件是否存在
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		common.CreateResponse(c, common.FileIsNotExists, "file not found")
		return
	}

	// 检查目标是否已存在
	if _, err := os.Stat(newPath); err == nil {
		common.CreateResponse(c, common.ErrorCode, "target name already exists")
		return
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		ylog.Errorf("RenameFile", "重命名失败: %v", err)
		common.CreateResponse(c, common.ErrorCode, "Failed to rename")
		return
	}

	ylog.Infof("RenameFile", "renamed: %s -> %s", oldPath, newPath)
	common.CreateResponse(c, common.SuccessCode, "rename successful")
}

// CreateUploadToken 创建上传凭证
func CreateUploadToken(c *gin.Context) {
	type MessageData struct {
		AccessKey string `json:"appkey" binding:"required"`
		SecretKey string `json:"appsecret" binding:"required"`
	}

	var data MessageData
	if err := c.ShouldBindJSON(&data); err != nil {
		common.BadRequest(c, err.Error())
		return
	}

	p := storage.UploadPolicy{}
	mac := auth.New(data.AccessKey, data.SecretKey)
	token := p.UploadToken(mac)

	common.CreateResponse(c, common.SuccessCode, token)
}
