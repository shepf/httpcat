package v1

import (
	"crypto/md5"
	"encoding/hex"
	"httpcat/internal/common"
	"httpcat/internal/common/ylog"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func DownloadFile(c *gin.Context) {

	// 从请求参数获取文件名
	fileName := c.Query("filename")

	path := filepath.Join(common.DownloadDir, fileName)
	//打印
	ylog.Infof("downloadFile", "download file from: %s", path)

	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		c.AbortWithStatus(404)
		ylog.Errorf("downloadFile", "打开文件失败,文件不存在", err)
		common.CreateResponse(c, common.FileIsNotExists, nil)
		return
	}
	// 获取文件信息
	fileInfo, _ := file.Stat()
	fileSize := int(fileInfo.Size())

	// 设置HEADER信息
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Length", strconv.FormatInt(int64(fileSize), 10))

	// 流式传输文件数据，并计算 MD5 值 相当于边下载边计算 MD5 值
	// 创建 MD5 哈希计算器
	hash := md5.New()

	// 流式传输文件数据，并计算 MD5 值
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if n == 0 {
			break
		}
		if err != nil && err != io.EOF {
			ylog.Errorf("DownloadFile", "读取文件失败", err)
			break
		}

		// 写入下载响应
		c.Writer.Write(buf[:n])

		// 计算 MD5 值
		hash.Write(buf[:n])
	}
	fileMD5 := hex.EncodeToString(hash.Sum(nil))

	// 获取其他文件信息
	createdTime := fileInfo.ModTime().Format("2006-01-02 15:04:05")
	modifiedTime := fileInfo.ModTime().Format("2006-01-02 15:04:05")

	// 异步记录下载日志
	go func() {
		recordDownloadLog(fileName, c.Request.RemoteAddr, fileSize, createdTime, modifiedTime, fileMD5)
	}()

}

func recordDownloadLog(fileName string, ip string, fileSize int, createdTime string, modifiedTime string, fileMD5 string) {
	dbPath := common.SqliteDBPath
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		ylog.Errorf("recordDownloadLog", "打开数据库失败", err)
		return
	}
	db.Debug()

	//// 自动创建表，服务初始化的时候，已经创建了表，这里不需要再创建
	//err = db.AutoMigrate(&common.DownloadLogModel{})
	//if err != nil {
	//	ylog.Errorf("recordDownloadLog", "创建表失败", err)
	//	return
	//}

	log := common.DownloadLogModel{
		IP:           ip,
		AppKey:       "", // 根据实际情况设置 AppKey
		DownloadTime: time.Now().Format("2006-01-02 15:04:05"),
		FileName:     fileName,
		FileSize:     strconv.Itoa(fileSize),
		FileMD5:      fileMD5,
		CreatedTime:  createdTime,
		ModifiedTime: modifiedTime,
	}

	err = db.Create(&log).Error
	if err != nil {
		ylog.Errorf("recordDownloadLog", "记录下载日志失败", err)
	}
}
