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
	//定义根路径路由,显示首页
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// 使用两个路由组,一个处理静态文件,一个处理 API
	//在这个静态文件组中,以`common.StaticDir`作为静态文件根目录,并且映射到路由`/static`下。这样访问:
	static := router.Group("/static")
	static.StaticFS("/", http.Dir(common.StaticDir))
}

func RegisterRouter(r *gin.Engine) {

	registerForFrontEnd(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"commit":  common.Commit,
			"build":   common.Build,
			"version": common.Version,
			"ci":      common.CI,
		})
	})

	r.Use(midware.Metrics())

	var (
		apiv1Group *gin.RouterGroup
		//apiv2Group *gin.RouterGroup
	)
	r.Use(Cors())

	apiv1Group = r.Group("/api/v1")
	{
		//apiv1Group.Use(midware.TokenAuth())
		//apiv1Group.Use(midware.RBACAuth())

		//用户操作相关接口
		userRouter := apiv1Group.Group("/user")
		{
			userRouter.POST("/login/account", v1.UserLogin)
			//	userRouter.GET("/logout", v1.UserLoginout)
			//	userRouter.POST("/del", v1.DelUser)
			//	userRouter.GET("/info", v1.UserInfo)
			//	userRouter.POST("/update", v1.UpdateUser)
			//	userRouter.POST("/resetPassword", v1.ResetPassword)
			//	userRouter.POST("/checkUser", v1.CheckPassword)
			userRouter.POST("/createUploadToken", createUploadToken)
			userRouter.POST("/checkUploadToken", checkUploadToken)

		}

		if common.FileEnable {
			ylog.Infof("[ROUTE]", "httpcat 开启 文件上传下载功能")
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
				fileRouter.GET("/download", downloadFile)
				// 获取目录文件列表
				fileRouter.GET("/listFiles", listFiles)
				// 获取某个文件的信息
				fileRouter.GET("/fileInfo", fileInfo)
			}
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
