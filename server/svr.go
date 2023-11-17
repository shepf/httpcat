package server

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gin_web_demo/server/common"
	"gin_web_demo/server/common/ylog"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

func RunAPIServer(port int, enableSSL, enableAuth bool, certFile, keyFile string) {

	//生成一个 Engine，这是 gin 的核心，默认带有 Logger 和 Recovery 两个中间件
	router := gin.Default()
	RegisterRouter(router)

	// 创建http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
		//ReadTimeout：服务器在读取客户端请求时，等待的最大时间。
		//如果设置为X分钟，那么如果服务器在X分钟内没有读取到完整的客户端请求，那么就会返回一个超时错误。
		ReadTimeout: time.Duration(common.HttpReadTimeout) * time.Second,
		// 服务器在写回应应答时，等待的最大时间。如果设置为X分钟，那么如果服务器在X分钟内没有写完应答，那么就会返回一个超时错误。
		WriteTimeout: time.Duration(common.HttpWriteTimeout) * time.Second,
		// 一个连接在空闲状态下（即没有任何数据传输），可以存在的最长时间。
		IdleTimeout: time.Duration(common.HttpIdleTimeout) * time.Second,
	}

	var err error
	ylog.Infof("RunServer", "####HTTP_LISTEN_ON:%d", port)
	if enableSSL {
		// 用ListenAndServeTLS替代router.RunTLS
		err = srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		// 用srv.ListenAndServe()替代router.Run
		err = srv.ListenAndServe()
	}
	if err != nil {
		ylog.Errorf("RunServer", "####http run error: %v", err)
	}

}

func getDirConf(c *gin.Context) {

	dirConf := make(map[string]string)
	dirConf["UploadDir"] = common.UploadDir
	dirConf["DownloadDir"] = common.DownloadDir
	dirConf["StaticDir"] = common.StaticDir

	common.CreateResponse(c, common.SuccessCode, dirConf)

}

func uploadFile(c *gin.Context) {
	// FormFile方法会读取参数“upload”后面的文件名，返回值是一个File指针，和一个FileHeader指针，和一个err错误。
	file, header, err := c.Request.FormFile("f1")
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	// header调用Filename方法，就可以得到文件名
	filename := header.Filename
	fmt.Println(file, err, filename)

	filePath := common.UploadDir + filename
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

	common.CreateResponse(c, common.SuccessCode, "upload successful!")
}

func downloadFile(c *gin.Context) {

	// 从请求参数获取文件名
	fileName := c.Query("filename")

	path := filepath.Join(common.DownloadDir, fileName)
	//打印
	ylog.Infof("downloadFile", "download file from: %s", path)

	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	// 获取文件信息
	fileInfo, _ := file.Stat()
	fileSize := int(fileInfo.Size())

	// 设置HEADER信息
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Length", strconv.FormatInt(int64(fileSize), 10))

	// 流式传输文件数据
	c.Stream(func(w io.Writer) bool {
		buf := make([]byte, 1024)
		for {
			n, _ := file.Read(buf)
			if n == 0 {
				break
			}
			w.Write(buf[:n])
		}
		return false
	})

}

// 获取目录文件列表
func listFiles(c *gin.Context) {

	dirPath := c.Query("dir")

	// 检查目录路径
	//if !strings.HasPrefix(dirPath, common.UploadDir) {
	//	c.AbortWithStatus(403)
	//	return
	//}

	dirPath = common.DownloadDir + dirPath
	ylog.Infof("listFiles func:", "dirPath:%s", dirPath)

	// 读取目录
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			ylog.Errorf("目录不存在", "err:%v", err)
			common.CreateResponse(c, common.DirISNotExists, "Directory does not exist")
		} else {
			ylog.Errorf("读取目录失败", "err:%v", err)
			common.CreateResponse(c, common.ReadDirFailed, "Failed to read the directory")
		}
		c.AbortWithStatus(500)
		return
	}

	// 按照文件时间倒序排列
	sort.SliceStable(files, func(i, j int) bool {
		return files[j].ModTime().Before(files[i].ModTime())
	})

	// 构建返回结果
	var fileList []map[string]interface{}
	for _, fileInfo := range files {
		fileEntry := make(map[string]interface{})
		fileEntry["FileName"] = fileInfo.Name()
		fileEntry["LastModified"] = fileInfo.ModTime().Format("2006-01-02 15:04:05")
		fileEntry["Size"] = formatSize(fileInfo.Size())
		fileList = append(fileList, fileEntry)
	}

	// 返回文件列表
	common.CreateResponse(c, common.SuccessCode, fileList)

}

func formatSize(size int64) string {
	const (
		B = 1 << (10 * iota)
		KB
		MB
		GB
		TB
		PB
	)

	switch {
	case size >= PB:
		return fmt.Sprintf("%.2f PB", float64(size)/PB)
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/TB)
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	}
	return fmt.Sprintf("%d B", size)
}

func fileInfo(c *gin.Context) {
	fileName := c.Query("name")

	// 检查文件路径
	filePath := common.DownloadDir + fileName
	ylog.Infof("fileInfo func:", "filePath:%s", filePath)

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			ylog.Errorf("文件不存在", "err:%v", err)
			common.CreateResponse(c, common.ErrorCode, "File does not exist")
		} else {
			ylog.Errorf("获取文件信息失败", "err:%v", err)
			common.CreateResponse(c, common.ErrorCode, "Failed to get file information")
		}
		c.AbortWithStatus(500)
		return
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		ylog.Errorf("打开文件失败", "err:%v", err)
		common.CreateResponse(c, common.ErrorCode, "Failed to open the file")
		c.AbortWithStatus(500)
		return
	}
	defer file.Close()

	// 创建一个MD5 hash
	hash := md5.New()

	// 将文件的内容复制到hash中
	if _, err := io.Copy(hash, file); err != nil {
		c.AbortWithStatus(500)
		return
	}

	// 获取MD5 hash的值
	md5Hash := hex.EncodeToString(hash.Sum(nil))

	// 构建返回结果
	fileEntry := make(map[string]interface{})
	fileEntry["FileName"] = fileInfo.Name()
	fileEntry["LastModified"] = fileInfo.ModTime().Format("2006-01-02 15:04:05")
	fileEntry["Size"] = formatSize(fileInfo.Size())
	fileEntry["MD5"] = md5Hash

	// 返回文件信息
	common.CreateResponse(c, common.SuccessCode, fileEntry)
}
