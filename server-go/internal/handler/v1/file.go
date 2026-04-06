package v1

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"httpcat/internal/common"
	"httpcat/internal/common/ylog"
)

func DownloadFile(c *gin.Context) {
	fileName := c.Query("filename")
	path, err := common.ResolvePathWithinBase(common.GetDownloadDir(), fileName)
	if err != nil {
		common.BadRequest(c, "invalid filename")
		return
	}

	ylog.Infof("downloadFile", "download file from: %s", path)

	file, err := os.Open(path)
	if err != nil {
		c.AbortWithStatus(404)
		ylog.Errorf("downloadFile", "打开文件失败,文件不存在: %v", err)
		common.CreateResponse(c, common.FileIsNotExists, nil)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil || fileInfo.IsDir() {
		common.BadRequest(c, "invalid filename")
		return
	}
	fileSize := int(fileInfo.Size())

	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(fileInfo.Name()))
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Length", strconv.FormatInt(int64(fileSize), 10))

	hash := md5.New()
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if n == 0 {
			break
		}
		if err != nil && err != io.EOF {
			ylog.Errorf("DownloadFile", "读取文件失败: %v", err)
			break
		}

		_, _ = c.Writer.Write(buf[:n])
		_, _ = hash.Write(buf[:n])
	}
	fileMD5 := hex.EncodeToString(hash.Sum(nil))

	createdTime := fileInfo.ModTime().Format("2006-01-02 15:04:05")
	modifiedTime := fileInfo.ModTime().Format("2006-01-02 15:04:05")

	go func() {
		recordDownloadLog(fileInfo.Name(), c.Request.RemoteAddr, fileSize, createdTime, modifiedTime, fileMD5)
	}()
}

func recordDownloadLog(fileName string, ip string, fileSize int, createdTime string, modifiedTime string, fileMD5 string) {
	db, err := common.GetDB()
	if err != nil {
		ylog.Errorf("recordDownloadLog", "获取数据库连接失败: %v", err)
		return
	}

	log := common.DownloadLogModel{
		IP:           ip,
		AppKey:       "",
		DownloadTime: time.Now().Format("2006-01-02 15:04:05"),
		FileName:     fileName,
		FileSize:     strconv.Itoa(fileSize),
		FileMD5:      fileMD5,
		CreatedTime:  createdTime,
		ModifiedTime: modifiedTime,
	}

	err = db.Create(&log).Error
	if err != nil {
		ylog.Errorf("recordDownloadLog", "记录下载日志失败: %v", err)
	}
}
