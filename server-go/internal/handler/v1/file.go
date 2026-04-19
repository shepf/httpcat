package v1

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"httpcat/internal/common"
	"httpcat/internal/common/ylog"
)

// DownloadFile 文件下载（v0.7.0：支持 HTTP Range，断点续传/拖动进度条）
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

	// 设置响应头（ServeContent 会自动处理 Content-Length/Content-Range/If-Modified-Since 等）
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(fileInfo.Name()))
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Accept-Ranges", "bytes")

	// 记录下载日志（异步；Range 请求可能多次触发，这里粗略记录一次即可）
	// 仅在非 Range 请求或 Range 起始为 0 时记录，避免重复刷日志
	rangeHeader := c.Request.Header.Get("Range")
	shouldLog := rangeHeader == "" || rangeHeader == "bytes=0-"

	// 计算 MD5（仅在需要记录日志时计算，避免大文件每次下载都算 MD5）
	fileMD5 := ""
	if shouldLog {
		if h, err := computeFileMD5(path); err == nil {
			fileMD5 = h
		}
	}

	// 使用 http.ServeContent 处理内容发送：
	// - 支持 Range 请求（断点续传）
	// - 支持 If-Modified-Since 304
	// - 自动处理 Content-Length
	http.ServeContent(c.Writer, c.Request, fileInfo.Name(), fileInfo.ModTime(), file)

	if shouldLog {
		size := int(fileInfo.Size())
		createdTime := fileInfo.ModTime().Format("2006-01-02 15:04:05")
		modifiedTime := createdTime
		ip := c.ClientIP()
		go recordDownloadLog(fileInfo.Name(), ip, size, createdTime, modifiedTime, fileMD5)
	}
}

// computeFileMD5 计算文件 MD5
func computeFileMD5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
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

