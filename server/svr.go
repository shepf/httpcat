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
	"strconv"
	"strings"
)

func RunAPIServer(port int, enableSSL, enableAuth bool, certFile, keyFile string) {

	//生成一个 Engine，这是 gin 的核心，默认带有 Logger 和 Recovery 两个中间件
	router := gin.Default()

	//定义根路径路由,显示首页
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// 使用两个路由组,一个处理静态文件,一个处理 API
	//在这个静态文件组中,以`common.StaticDir`作为静态文件根目录,并且映射到路由`/static`下。这样访问:
	static := router.Group("/static")
	static.StaticFS("/", http.Dir(common.StaticDir))

	api := router.Group("/api")
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	api.POST("/upload", uploadFile)
	//使用实现 API 方式进行文件下载,而不是直接通过 StaticFS 暴露文件：
	//原因是:
	//1. StaticFS 更适合提供静态资源文件的访问,这些文件通常对所有用户都是公开的,不需要鉴权。
	//2. 对于需要权限控制的文件下载,实现 API 方式更合适,可以方便地在代码中添加鉴权逻辑。
	api.GET("/download", downloadFile)
	// 获取目录文件列表
	api.GET("/listFiles", listFiles)

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

	path := common.DownloadDir + fileName
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
	if !strings.HasPrefix(dirPath, common.UploadDir) {
		c.AbortWithStatus(403)
		return
	}

	// 读取目录
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	// 返回文件名和文件信息
	c.JSON(200, files)

}
