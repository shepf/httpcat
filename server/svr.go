package server

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gin_web_demo/server/common"
	"gin_web_demo/server/common/utils"
	"gin_web_demo/server/common/ylog"
	"gin_web_demo/server/p2p"
	"gin_web_demo/server/storage"
	"gin_web_demo/server/storage/auth"
	"github.com/gin-gonic/gin"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
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

var ps *pubsub.PubSub

// 使用一个全局变量 subscribedTopics，它是一个映射（map），用于保存已加入的主题及其对应的 *pubsub.Topic 实例
var subscribedTopics = make(map[string]*pubsub.Topic)

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
	ctx := context.Background()

	enableP2P := common.P2pEnable
	if enableP2P {
		go runP2PServer(ctx, router)
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

func runP2PServer(ctx context.Context, router *gin.Engine) {

	// To construct a simple host with all the default settings, just use `New`
	ip := common.P2pListenIP
	port := common.P2pListenPort
	listenAddr := fmt.Sprintf("/ip4/%s/tcp/%d", ip, port)

	fmt.Printf("[*] Listening on: %s with port: %d\n", ip, port)
	fmt.Println("p2p listenAddr:", listenAddr)

	h, err := libp2p.New(
		libp2p.ListenAddrStrings(
			listenAddr, // "/ip4/0.0.0.0/tcp/9000", // regular tcp connections
			//"/ip4/0.0.0.0/udp/9000/quic", // a UDP endpoint for the QUIC transport
		),
	)

	if err != nil {
		panic(err)
	}
	defer h.Close()

	fmt.Printf("\033[32mHello World, my p2p hosts ID is %s\033[0m\n", h.ID())

	// 节点发现
	go discoverPeers(ctx, h)

	// PubSub
	ps, _ = pubsub.NewGossipSub(ctx, h)
	ylog.Infof("runP2PServer", "join topic: %v", common.TopicName)
	topic, _ := ps.Join(common.TopicName)
	subscription, _ := topic.Subscribe()
	ylog.Infof("runP2PServer", "subscribed topic: %v", common.TopicName)

	// 在 libp2p 的 pubsub 模型中，不支持直接获取当前节点已订阅的主题列表
	// 自己维护一个主题列表。每当节点加入一个主题时，将其添加到该列表中
	// 将主题添加到已订阅的主题列表
	// 将主题添加到已订阅的主题列表
	subscribedTopics[common.TopicName] = topic

	go func() {
		for {
			msg, err := subscription.Next(ctx)
			if err != nil {
				// handle error
				break
			}
			fmt.Printf("Received message from %s: %s\n", msg.GetFrom(), string(msg.GetData()))

		}
	}()

	// 等待上下文取消信号
	<-ctx.Done()
	fmt.Println("P2P server stopped")
}

func publishMessage(c *gin.Context, topicName string, message string) {

	// 检查是否已经加入主题
	topic, exists := subscribedTopics[topicName]
	if !exists {
		// 加入主题
		var err error
		// 经过测试，在 libp2p 的 pubsub 模型中，当您重复调用 Join 方法加入相同的主题时，会报错：topic already exists
		topic, err = ps.Join(topicName)
		if err != nil {
			// 处理错误
			ylog.Errorf("publishMessage", "Failed to join the topic, err:%v", err)
			common.CreateResponse(c, common.ErrorCode, "Failed to join the topic")
			return
		}

		// 将主题添加到已订阅的主题列表
		subscribedTopics[topicName] = topic
	}

	// 发布消息到主题
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := topic.Publish(ctx, []byte(message))
	if err != nil {
		// handle error
		common.CreateResponse(c, common.ErrorCode, "Failed to publish the message")
		return
	}

	common.CreateResponse(c, common.SuccessCode, "Message sent to topic: "+topicName)

}

func getSubscribedTopics(c *gin.Context) {

	common.CreateResponse(c, common.SuccessCode, subscribedTopics)
	return
}

func discoverPeers(ctx context.Context, h host.Host) {
	if common.EnableMdns {
		fmt.Printf("Host ID is %s. Enabling MDNS for discovering nodes!\n", h.ID())

		peerChan := p2p.InitMDNS(h, common.RendezvousString)

		// 维护一个连接的节点列表
		connectedPeers := map[peer.ID]bool{}

		// Look for others who have announced and attempt to connect to them
		for {
			peer := <-peerChan // will block until we discover a peer
			if peer.ID == h.ID() {
				continue // No self connection
			}

			fmt.Println("Found peer:")
			fmt.Println("ID:", peer.ID)

			// 获取节点的地址
			if _, ok := connectedPeers[peer.ID]; !ok {
				if err := h.Connect(ctx, peer); err != nil {
					fmt.Println("Connection failed:", err)
					continue
				}

				// 添加到已连接节点列表
				connectedPeers[peer.ID] = true

				// 打印已连接的节点
				fmt.Println("Connected to:", peer.ID)
				fmt.Println("Connected peers:")
				for connectedPeer := range connectedPeers {
					fmt.Println("- ", connectedPeer)
				}
			}
		}
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

	// 判断是否开启了UploadToken校验
	if common.EnableUploadToken {
		uploadToken := c.Request.Header.Get("UploadToken")
		if uploadToken == "" {
			common.CreateResponse(c, common.ErrorCode, "UploadToken is empty")
			return
		}
		// 校验UploadToken
		accessKey := common.AppKey
		secretKey := common.AppSecret
		mac := auth.New(accessKey, secretKey)

		if !mac.VerifyUploadToken(uploadToken) {
			common.CreateResponse(c, common.ErrorCode, "UploadToken is invalid")
			return
		}
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

	fmt.Println("PersistentNotifyURL:", common.PersistentNotifyURL)
	// 上传成功后，发送通知
	if common.PersistentNotifyURL != "" {

		ylog.Infof("uploadFile", "send notify to: %s", common.PersistentNotifyURL)
		// 构建 Markdown 通知内容
		ip := c.ClientIP()
		uploadTime := time.Now().Format("2006-01-02 15:04:05")
		filename := header.Filename
		// 获取文件信息
		fileInfo, _ := os.Stat(filePath)
		fileSize := formatSize(fileInfo.Size())
		fileMD5, _ := utils.CalculateMD5(filePath)

		markdownContent := fmt.Sprintf(`>有文件上传归档,上传信息：
			- IP地址：%s
			- 上传时间：%s
			- 文件名：%s
			- 文件大小：%s
			- 文件MD5：%s`, ip, uploadTime, filename, fileSize, fileMD5)
		ylog.Infof("uploadFile", "markdownContent:%s", markdownContent)

		go utils.SendNotify(common.PersistentNotifyURL, markdownContent)

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

func sendP2pMessage(c *gin.Context) {
	// 定义结构体用于解析 JSON 数据
	type MessageData struct {
		Topic   string `json:"topic"`
		Message string `json:"message"`
	}

	var data MessageData

	// 解析 JSON 数据到结构体
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON data",
		})
		return
	}

	// 获取主题和消息内容
	topic := data.Topic
	message := data.Message

	// Publish a message to the topic
	publishMessage(c, topic, message)

}

func createUploadToken(c *gin.Context) {
	// 定义结构体用于解析 JSON 数据
	type MessageData struct {
		AccessKey string `json:"accessKey"`
		SecretKey string `json:"secretKey"`
	}

	var data MessageData

	// 解析 JSON 数据到结构体
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON data",
		})
		return
	}

	// 获取主题和消息内容
	accessKey := data.AccessKey
	secretKey := data.SecretKey

	p := storage.UploadPolicy{}
	mac := auth.New(accessKey, secretKey)
	token := p.UploadToken(mac)

	common.CreateResponse(c, common.SuccessCode, token)
}

func checkUploadToken(c *gin.Context) {
	// 获取请求头中的UploadToken
	uploadToken := c.Request.Header.Get("UploadToken")

	// 校验UploadToken
	accessKey := common.AppKey
	secretKey := common.AppSecret
	mac := auth.New(accessKey, secretKey)

	if !mac.VerifyUploadToken(uploadToken) {
		common.CreateResponse(c, common.ErrorCode, "UploadToken is invalid")
		return
	}

	common.CreateResponse(c, common.SuccessCode, "UploadToken is valid")

}
