package server

import (
	"fmt"
	"gin_web_demo/server/common"
	"gin_web_demo/server/common/ylog"
	v1 "gin_web_demo/server/handler/v1"
	"gin_web_demo/server/midware"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func registerForFrontEnd(router *gin.Engine) {

	StaticDir := common.StaticDir
	//定义根路径路由,显示首页
	router.GET("/", func(c *gin.Context) {
		c.File(StaticDir + "/index.html")
	})

	fmt.Println("StaticDir:", StaticDir)
	// 处理静态文件
	router.StaticFS("/static", http.Dir(StaticDir))

	// 定义通配符路由，将所有未匹配到的路由重定向到 index.html
	router.NoRoute(func(c *gin.Context) {
		c.File(StaticDir + "/index.html")
	})

}

func RegisterRouter(r *gin.Engine) {

	registerForFrontEnd(r)

	r.Use(midware.Metrics())

	var (
		apiv1Group *gin.RouterGroup
		//apiv2Group *gin.RouterGroup
	)
	r.Use(Cors())

	apiv1Group = r.Group("/api/v1")
	{
		apiv1Group.Use(midware.TokenAuth())
		//apiv1Group.Use(midware.RBACAuth())

		confRouter := apiv1Group.Group("/conf")
		{
			confRouter.GET("/getVersion", v1.GetVersion)
			confRouter.GET("/getConf", v1.GetConfInfo)

		}

		//用户操作相关接口
		userRouter := apiv1Group.Group("/user")
		{
			userRouter.POST("/login/account", v1.UserLogin)
			userRouter.GET("/currentUser", v1.UserInfo)
			userRouter.POST("/login/outLogin", v1.UserLoginout)
			userRouter.POST("/changePasswd", v1.ChangePasswd)
			//修改用户信息
			//上传用户头像
			userRouter.POST("/uploadAvatar", v1.UploadAvatar)

			//	userRouter.POST("/del", v1.DelUser)
			//	userRouter.POST("/update", v1.UpdateUser)
			//	userRouter.POST("/resetPassword", v1.ResetPassword)
			//	userRouter.POST("/checkUser", v1.CheckPassword)
			userRouter.GET("/generateAppSecret", generateAppSecret)
			userRouter.POST("/saveUploadToken", saveUploadToken)
			userRouter.DELETE("/removeUploadToken", removeUploadToken)
			userRouter.GET("/uploadTokenLists", getUploadTokenLists)
			userRouter.POST("/createUploadToken", createUploadToken)
			userRouter.POST("/checkUploadToken", checkUploadToken)

			// 统计信息
			// 数据概览 Data Overview
			userRouter.GET("/dataOverview", dataOverview)
			userRouter.GET("/getUploadAvailableSpace", getUploadAvailableSpace)

		}

		// 统计信息
		// 数据概览 Data Overview
		statisticsRouter := apiv1Group.Group("/statistics")
		{
			statisticsRouter.GET("/getUploadStatistics", v1.GetUploadStatistics)
			statisticsRouter.GET("/getDownloadStatistics", v1.GetDownloadStatistics)
		}

		// 文件操作相关接口
		fileRouter := apiv1Group.Group("/file")
		{
			//获取配置文件中的上传下载目录配置
			fileRouter.GET("/getDirConf", getDirConf)
			fileRouter.POST("/upload", uploadFile)
			//使用实现 API 方式进行文件下载,而不是直接通过 StaticFS 暴露文件：
			//原因是:
			//1. StaticFS 更适合提供静态资源文件的访问,这些文件通常对所有用户都是公开的,不需要鉴权。
			//2. 对于需要权限控制的文件下载,实现 API 方式更合适,可以方便地在代码中添加鉴权逻辑。
			fileRouter.GET("/download", v1.DownloadFile)
			// 获取目录文件列表
			fileRouter.GET("/listFiles", listFiles)
			// 获取某个文件的信息
			fileRouter.GET("/getFileInfo", getFileInfo)

			//获取上传文件历史记录
			fileRouter.GET("/uploadHistoryLogs", uploadHistoryLogs)
			// 删除上传历史记录
			fileRouter.DELETE("/uploadHistoryLogs", deleteHistoryLogs)

		}

		imageManageRouter := apiv1Group.Group("/imageManage")
		{
			imageManageRouter.POST("/upload", v1.UploadImage)
			//文件改名
			imageManageRouter.POST("/rename", v1.RenameImage)
			//图片文件删除
			imageManageRouter.DELETE("/delete", v1.DeleteImage)
			// 下载图片
			imageManageRouter.GET("/download", v1.DownloadImage)
			// 分页获取图片缩略图
			imageManageRouter.GET("/listThumbImages", v1.GetThumbnails)

		}

		if common.P2pEnable {
			ylog.Infof("[ROUTE]", "httpcat 开启 P2P功能")
			// 文件操作相关接口
			p2pRouter := apiv1Group.Group("/p2p")
			{
				p2pRouter.POST("/send_message", sendP2pMessage)
				p2pRouter.POST("/get_subscribed_topics", getSubscribedTopics)
			}
		}

	}

}

// 这段代码的功能是处理跨域请求的设置，允许不同域的客户端进行访问，并设置相应的头部信息以满足跨域请求的要求。
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {

			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar,Content-Disposition, token")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}

		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "")
		}
		c.Next()
	}
}
