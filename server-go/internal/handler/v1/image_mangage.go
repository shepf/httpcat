package v1

import (
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"httpcat/internal/common"
	"httpcat/internal/common/utils"
	"httpcat/internal/common/ylog"
	"httpcat/internal/midware"
	"httpcat/internal/models"
	"httpcat/internal/storage/auth"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func UploadImage(c *gin.Context) {

	if !common.FileUploadEnable {
		common.CreateResponse(c, common.ErrorCode, "File service is not enabled")
		return
	}

	jwtAuth := false
	jwtUsername := ""
	if authHeader := c.Request.Header.Get("Authorization"); authHeader != "" {
		tokenStr := authHeader
		if len(tokenStr) > 7 && strings.ToUpper(tokenStr[0:7]) == "BEARER " {
			tokenStr = tokenStr[7:]
		}
		if tokenStr != "" {
			if claims, err := midware.VerifyToken(tokenStr, []byte(common.JwtSecret)); err == nil {
				jwtAuth = true
				if username, ok := (*claims)["username"].(string); ok {
					jwtUsername = username
				}
			}
		}
	}

	if !jwtAuth && common.EnableUploadToken {
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
		appkey := parts[0]

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
	} else if !jwtAuth {
		common.Unauthorized(c, "Authentication required")
		return
	}

	if jwtAuth && common.MustChangePassword(common.GetUser(jwtUsername)) {
		c.JSON(http.StatusForbidden, gin.H{
			"errorCode": common.PasswordNeedChanged,
			"msg":       common.ErrorDescriptions[common.PasswordNeedChanged],
		})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		common.BadRequest(c, "Bad request,check your file~")
		return
	}
	defer file.Close()

	filename, err := common.NormalizeSafeFileName(header.Filename)
	if err != nil {
		common.BadRequest(c, "invalid filename")
		return
	}
	fmt.Println(file, err, filename)

	imagesDir, err := common.ResolvePathWithinBase(common.GetUploadDir(), "images")
	if err != nil {
		common.CreateResponse(c, common.ErrorCode, "Failed to resolve images directory")
		return
	}
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		ylog.Errorf("uploadFile", "创建目录失败: %v", err)
		common.CreateResponse(c, common.ErrorCode, "Failed to prepare images directory")
		return
	}

	filePath, err := common.ResolvePathWithinBase(imagesDir, filename)
	if err != nil {
		common.BadRequest(c, "invalid filename")
		return
	}

	if _, err := os.Stat(filePath); err == nil {
		common.BadRequest(c, "File already exists")
		return
	}

	ylog.Infof("uploadFile", "upload file to: %s", filePath)
	out, err := os.Create(filePath)
	if err != nil {
		ylog.Errorf("uploadImage", "创建文件失败: %v", err)
		common.CreateResponse(c, common.ErrorCode, fmt.Sprintf("Failed to create file: %v", err))
		return
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		ylog.Errorf("uploadImage", "写入文件失败: %v", err)
		common.CreateResponse(c, common.ErrorCode, fmt.Sprintf("Failed to save file: %v", err))
		return
	}

	thumbFilePath, err := common.ResolvePathWithinBase(imagesDir, "thumb_"+filename)
	if err != nil {
		os.Remove(filePath)
		common.CreateResponse(c, common.ErrorCode, "Invalid thumbnail path")
		return
	}
	thumbImage, err := imaging.Open(filePath)
	if err != nil {
		ylog.Errorf("uploadImage", "解析图片失败（文件可能不是有效图片）: %v", err)
		os.Remove(filePath)
		common.CreateResponse(c, common.ErrorCode, fmt.Sprintf("Invalid image file, failed to parse: %v", err))
		return
	}
	thumbImage = imaging.Resize(thumbImage, common.ThumbWidth, common.ThumbHeight, imaging.Lanczos)
	if err = imaging.Save(thumbImage, thumbFilePath); err != nil {
		ylog.Errorf("uploadImage", "保存缩略图失败: %v", err)
		common.CreateResponse(c, common.ErrorCode, fmt.Sprintf("Failed to generate thumbnail: %v", err))
		return
	}

	Ip := c.ClientIP()
	uploadTime := time.Now().Format("2006-01-02 15:04:05")
	fileMD5, _ := utils.CalculateMD5(filePath)
	fileUUID := uuid.NewV4().String()

	if common.EnableSqlite {
		ylog.Infof("uploadFile", "sqliteInsert enable")

		db, err := common.GetDB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		image := models.UploadImageModel{
			FileUUID:      fileUUID,
			Size:          header.Size,
			FileName:      filename,
			FilePath:      filePath,
			ThumbFilePath: thumbFilePath,
			FileMD5:       fileMD5,
			DownloadCount: 0,
			Sort:          1000,
			UploadTime:    uploadTime,
			UploadIP:      Ip,
			UploadUser:    "admin",
			Status:        "done",
		}

		db.Create(&image)
	}

	c.JSON(http.StatusOK, gin.H{
		"errorCode": common.SuccessCode,
		"msg":       "upload successful!",
		"data": map[string]interface{}{
			"fileUUID":    fileUUID,
			"name":        filename,
			"status":      "done",
			"url":         "/api/v1/imageManage/download?filename=" + filename,
			"thumbUrl":    "/api/v1/imageManage/download?filename=thumb_" + filename,
			"description": "",
		},
	})

}

func RenameImage(c *gin.Context) {
	filename, err := common.NormalizeSafeFileName(c.PostForm("filename"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "文件名不合法"})
		return
	}
	newName, err := common.NormalizeSafeFileName(c.PostForm("newName"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "新文件名不合法"})
		return
	}

	imagesDir, err := common.ResolvePathWithinBase(common.GetUploadDir(), "images")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "图片目录无效"})
		return
	}
	filePath, err := common.ResolvePathWithinBase(imagesDir, filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "文件名不合法"})
		return
	}
	newFilePath, err := common.ResolvePathWithinBase(imagesDir, newName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "新文件名不合法"})
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "文件不存在"})
		return
	}

	err = os.Rename(filePath, newFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "重命名失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "图片名字修改成功"})
}

func DeleteImage(c *gin.Context) {
	filename, err := common.NormalizeSafeFileName(c.Query("filename"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "文件名不合法"})
		return
	}

	db, err := common.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库连接失败"})
		return
	}

	var image models.UploadImageModel
	err = db.Where("file_name = ?", filename).First(&image).Error
	if err == nil {
		db.Unscoped().Delete(&image)
	}

	imagesDir, err := common.ResolvePathWithinBase(common.GetUploadDir(), "images")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "图片目录无效"})
		return
	}
	filePath, err := common.ResolvePathWithinBase(imagesDir, filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "文件名不合法"})
		return
	}

	err = os.Remove(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "图片删除成功"})
}

func ClearImage(c *gin.Context) {

	db, err := common.GetDB()
	if err != nil {
		ylog.Errorf("clearImage", "数据库连接失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "数据库连接失败"})
		return
	}

	err = db.Exec("DELETE FROM t_upload_image").Error
	if err != nil {
		ylog.Errorf("clearImage", "清空数据库记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "清空照片失败"})
		return
	}

	dirPath, err := common.ResolvePathWithinBase(common.GetUploadDir(), "images")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "图片目录无效"})
		return
	}
	err = os.RemoveAll(dirPath)
	if err != nil {
		ylog.Errorf("clearImage", "清空照片失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "清空照片失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "照片清空成功"})
}

func DownloadImage(c *gin.Context) {
	filename, err := common.NormalizeSafeFileName(c.Query("filename"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid filename")
		return
	}
	imagesDir, err := common.ResolvePathWithinBase(common.GetUploadDir(), "images")
	if err != nil {
		c.String(http.StatusInternalServerError, "Invalid image directory")
		return
	}
	filePath, err := common.ResolvePathWithinBase(imagesDir, filename)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid filename")
		return
	}

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			c.String(http.StatusNotFound, "File not found")
			return
		}
		c.String(http.StatusInternalServerError, "Failed to read file")
		return
	}
	if info.IsDir() {
		c.String(http.StatusBadRequest, "Invalid filename")
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+strconv.Quote(info.Name()))
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
}

// 分页信息结构体
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"` //表示总页数，即符合查询条件的数据总共可以分成多少页。它是根据总记录数（或总项数）和每页显示的项数来计算得出的。
	TotalItems int64 `json:"totalItems"` // 表示总记录数或者总项数，即符合查询条件的所有数据项的数量
}

func GetThumbnails(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	search := c.DefaultQuery("search", "")

	// 将字符串转换为整数
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	imagesDir, err := common.ResolvePathWithinBase(common.GetUploadDir(), "images")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid image directory"})
		return
	}

	db, err := common.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var thumbnails []models.UploadImageModel
	var totalItems int64 // 声明 totalItems 变量并设置初始值为 0

	query := db.Model(&models.UploadImageModel{})
	// 支持按文件名模糊搜索
	if search != "" {
		query = query.Where("file_name LIKE ?", "%"+search+"%")
	}
	query.Count(&totalItems)

	// 计算总页数
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	query.Order("created_at desc, download_count desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&thumbnails)

	for i := range thumbnails {
		thumbnailPath := thumbnails[i].ThumbFilePath
		relThumbnailPath, err := filepath.Rel(imagesDir, thumbnailPath)
		if err != nil {
			ylog.Warnf("GetThumbnails", "skip thumbnail with invalid path %q: %v", thumbnailPath, err)
			continue
		}
		safeThumbnailPath, err := common.ResolvePathWithinBase(imagesDir, relThumbnailPath)
		if err != nil {
			ylog.Warnf("GetThumbnails", "skip thumbnail outside images dir %q: %v", thumbnailPath, err)
			continue
		}

		// 检查缩略图文件是否存在
		_, err = os.Stat(safeThumbnailPath)
		if os.IsNotExist(err) {
			// 缩略图不存在，跳过当前循环
			continue
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read thumbnail file"})
			return
		}

		// 读取缩略图文件
		fileBytes, err := os.ReadFile(safeThumbnailPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read thumbnail file"})
			return
		}

		// 将缩略图文件转换为 Base64 格式
		base64Image := base64.StdEncoding.EncodeToString(fileBytes)

		// 将 Base64 缩略图赋值给字段
		thumbnails[i].ThumbnailBase64 = base64Image
	}

	// 构建包含分页信息的响应数据
	response := struct {
		Pagination Pagination                `json:"pagination"`
		Data       []models.UploadImageModel `json:"data"`
	}{
		Pagination: Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			TotalItems: totalItems,
		},
		Data: thumbnails,
	}

	c.JSON(http.StatusOK, response)
}
