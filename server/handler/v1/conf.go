package v1

import (
	"gin_web_demo/server/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetVersion(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"commit":  common.Commit,
		"build":   common.Build,
		"version": common.Version,
		"ci":      common.CI,
	})

}
