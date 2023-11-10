package server

import (
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
)

func RunAPIServer(port int, enableSSL, enableAuth bool, certFile, keyFile string) {

	//生成一个 Engine，这是 gin 的核心，默认带有 Logger 和 Recovery 两个中间件
	router := gin.Default()
	RegisterRouter(router)

	var err error
	ylog.Infof("RunServer", "####HTTP_LISTEN_ON:%d", port)
	if enableSSL {
		err = router.RunTLS(fmt.Sprintf(":%d", port), certFile, keyFile)
	} else {
		err = router.Run(fmt.Sprintf(":%d", port))
	}
	if err != nil {
		ylog.Errorf("RunServer", "####http run error: %v", err)
	}

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
	log.Printf("upload file to: %s", filePath)
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

	c.String(http.StatusCreated, "upload successful \n")
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

	dirPath = common.UploadDir + dirPath

	// 读取目录
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
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
