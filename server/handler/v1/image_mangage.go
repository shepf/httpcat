package v1

import (
	"encoding/base64"
	"fmt"
	"gin_web_demo/server/common"
	"gin_web_demo/server/common/utils"
	"gin_web_demo/server/common/ylog"
	"gin_web_demo/server/models"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func UploadImage(c *gin.Context) {

	if !common.FileUploadEnable {
		common.CreateResponse(c, common.ErrorCode, "File service is not enabled")
		return
	}

	// FormFile方法会读取参数“upload”后面的文件名，返回值是一个File指针，和一个FileHeader指针，和一个err错误。
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		common.BadRequest(c, "Bad request,check your file~")
		return
	}

	// header调用Filename方法，就可以得到文件名
	filename := header.Filename
	fmt.Println(file, err, filename)

	// 保存文件到本地, 配置的上传目录加images目录
	filePath := filepath.Join(common.UploadDir, "images", filename)
	// 判断目录是否存在，如果不存在则创建
	imagesDir := filepath.Join(common.UploadDir, "images")
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		err := os.MkdirAll(imagesDir, 0755)
		if err != nil {
			ylog.Errorf("uploadFile", "创建目录失败", err)
			panic(err)
		}
	}

	// 判断文件是否存在
	_, err = os.Stat(filePath)
	if err == nil {
		// 文件已存在
		common.CreateResponse(c, common.ErrorCode, "File already exists")
		return
	}

	ylog.Infof("uploadFile", "upload file to: %s", filePath)
	// 创建一个文件，文件名为filename，这里的返回值out也是一个File指针
	out, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// 将file的内容拷贝到out
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}

	// 生成缩略图
	thumbFilePath := filepath.Join(common.UploadDir, "images", "thumb_"+filename)
	thumbImage, err := imaging.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	thumbImage = imaging.Resize(thumbImage, 250, 150, imaging.Lanczos) // 设置缩略图的宽度为 100
	err = imaging.Save(thumbImage, thumbFilePath)
	if err != nil {
		log.Fatal(err)
	}

	Ip := c.ClientIP()
	uploadTime := time.Now().Format("2006-01-02 15:04:05")
	// 获取文件信息
	fileMD5, _ := utils.CalculateMD5(filePath)
	fileUUID := uuid.NewV4().String()

	//// 是否sqlite记录
	if common.EnableSqlite {
		ylog.Infof("uploadFile", "sqliteInsert enable")

		dbPath := common.SqliteDBPath
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		db.Debug()

		// 保存图片信息到数据库
		image := models.UploadImageModel{
			FileUUID:      fileUUID,
			Size:          header.Size,
			FileName:      filename,
			FilePath:      filePath,
			ThumbFilePath: thumbFilePath,
			FileMD5:       fileMD5, // 计算文件的 MD5 值
			DownloadCount: 0,
			Sort:          1000,
			UploadTime:    uploadTime,
			UploadIP:      Ip,
			UploadUser:    "admin",
			Status:        "done",
		}

		db.Create(&image)

	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"errorCode": common.SuccessCode,
		"msg":       "upload successful!",
		"data": map[string]interface{}{
			"fileUUID": fileUUID,
			"name":     filename,
			"status":   "done",
			// 通常情况下，上传成功后前端需要再次请求图片的 URL 来展示图片。
			// 在上传成功后，后端会返回图片的 URL，前端可以使用这个 URL 来获取图片数据，并将其展示在页面上。
			"url":         "/api/v1/imageManage/download?filename=" + filename,
			"thumbUrl":    "/api/v1/imageManage/download?filename=thumb_" + filename,
			"description": "",
		},
	})

}

func RenameImage(c *gin.Context) {
	// 获取请求参数
	filename := c.PostForm("filename")
	newName := c.PostForm("newName")

	// 构建图片文件的完整路径
	filePath := filepath.Join(common.UploadDir, "images", filename)

	// 判断文件是否存在
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "文件不存在",
		})
		return
	}

	// 构建新的文件路径
	newFilePath := filepath.Join(common.UploadDir, "images", newName)

	// 重命名文件
	err = os.Rename(filePath, newFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "重命名失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "图片名字修改成功",
	})
}

func DeleteImage(c *gin.Context) {
	// 获取请求参数
	filename := c.Query("filename")

	// 从数据库中删除相应记录
	dbPath := common.SqliteDBPath
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "数据库连接失败",
		})
		return
	}

	// 删除数据库中的记录
	var image models.UploadImageModel
	err = db.Where("file_name = ?", filename).First(&image).Error
	//如果您使用的是 GORM，默认启用了软删除功能，即通过将 deleted_at 字段设置为非空来标记删除的记录。
	//这就是为什么在执行 db.Delete(&image) 后，实际上只是更新了 deleted_at 字段而不是真正删除记录。
	if err == nil {
		//db.Delete(&image)
		//如果您想要完全删除记录而不是软删除，可以使用 Unscoped() 方法来删除记录。
		db.Unscoped().Delete(&image)
	}

	// 构建图片文件的完整路径
	filePath := filepath.Join(common.UploadDir, "images", filename)

	//// 判断文件是否存在
	//_, err = os.Stat(filePath)
	//if os.IsNotExist(err) {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"message": "文件不存在",
	//	})
	//	return
	//}

	// 删除文件
	err = os.Remove(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "删除失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "图片删除成功",
	})
}

func ClearImage(c *gin.Context) {

	// 清空数据库中的记录
	dbPath := common.SqliteDBPath
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		ylog.Errorf("clearImage", "数据库连接失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "数据库连接失败",
		})
		return
	}

	// 清空数据库中的记录
	err = db.Exec("DELETE FROM t_upload_image").Error
	if err != nil {
		ylog.Errorf("clearImage", "清空数据库记录失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "清空照片失败",
		})
		return
	}

	// 删除图片文件夹下的所有文件
	dirPath := filepath.Join(common.UploadDir, "images")
	err = os.RemoveAll(dirPath)
	if err != nil {
		ylog.Errorf("clearImage", "清空照片失败", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "清空照片失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "照片清空成功",
	})
}

func DownloadImage(c *gin.Context) {
	filename := c.Query("filename")
	filePath := filepath.Join(common.UploadDir, "images", filename)

	// 检查文件是否存在
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			c.String(http.StatusNotFound, "File not found")
			return
		}
		c.String(http.StatusInternalServerError, "Failed to read file")
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
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

	// 将字符串转换为整数
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	dbPath := common.SqliteDBPath
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var thumbnails []models.UploadImageModel
	var totalItems int64 // 声明 totalItems 变量并设置初始值为 0
	db.Model(&models.UploadImageModel{}).Count(&totalItems)

	// 计算总页数
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&thumbnails)

	for i := range thumbnails {
		thumbnailPath := thumbnails[i].ThumbFilePath

		// 检查缩略图文件是否存在
		_, err := os.Stat(thumbnailPath)
		if os.IsNotExist(err) {
			// 缩略图不存在，跳过当前循环
			continue
		}

		// 读取缩略图文件
		fileBytes, err := ioutil.ReadFile(thumbnailPath)
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
