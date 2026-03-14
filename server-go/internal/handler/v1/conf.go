package v1

import (
	"fmt"
	"httpcat/internal/common"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/process"
	"time"
)

func GetVersion(c *gin.Context) {
	//在处理HTTP请求时获取进程的创建时间，它只会反映当前goroutine的创建时间，而不是整个应用程序的启动时间
	//pid := os.Getpid()
	//uptime, err := getProcessUptime(pid)
	//
	//if err != nil {
	//	log.Fatalf("Failed to get process uptime: %v", err)
	//}
	uptime := time.Since(common.StartTime)

	uptimeString := formatDuration(uptime)
	fmt.Printf("Process uptime: %s\n", uptimeString)

	common.CreateAntResponse(c, common.SuccessCode, gin.H{
		"commit":  common.Commit,
		"build":   common.Build,
		"version": common.Version,
		"ci":      common.CI,
		"uptime":  uptimeString,
	})

}

func GetConfInfo(c *gin.Context) {
	uploadDir := common.GetUploadDir()
	downloadDir := common.GetDownloadDir()
	webDir := common.StaticDir
	// 上传文件开关状态
	fileUploadEnable := common.FileUploadEnable

	// 获取工作目录
	workDir, err := os.Getwd()
	if err != nil {
		workDir = "-"
	}

	// 将相对路径转为绝对路径
	absUploadDir, _ := filepath.Abs(uploadDir)
	absDownloadDir, _ := filepath.Abs(downloadDir)
	absWebDir, _ := filepath.Abs(webDir)

	// 文件根目录：如果是默认值 "./" 或空，展示为项目工作目录
	fileBaseDir := common.FileBaseDir
	absFileBaseDir, _ := filepath.Abs(fileBaseDir)
	if fileBaseDir == "" || fileBaseDir == "./" || fileBaseDir == "." {
		fileBaseDir = workDir // 默认等于项目工作目录，直接展示绝对路径
	}

	common.CreateResponse(c, common.SuccessCode, gin.H{
		"uploadDir":        uploadDir,
		"downloadDir":      downloadDir,
		"webDir":           webDir,
		"fileUploadEnable": fileUploadEnable,
		"workDir":          workDir,
		"fileBaseDir":      fileBaseDir,
		"absFileBaseDir":   absFileBaseDir,
		"absUploadDir":     absUploadDir,
		"absDownloadDir":   absDownloadDir,
		"absWebDir":        absWebDir,
	})
}

func getProcessUptime(pid int) (time.Duration, error) {
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return 0, err
	}

	createTime, err := proc.CreateTime()
	if err != nil {
		return 0, err
	}

	uptime := time.Since(time.Unix(int64(createTime/1000), 0))

	return uptime, nil
}

func formatDuration(duration time.Duration) string {
	seconds := int(duration.Seconds())

	days := seconds / (60 * 60 * 24)
	seconds %= (60 * 60 * 24)

	hours := seconds / (60 * 60)
	seconds %= (60 * 60)

	minutes := seconds / 60
	seconds %= 60

	uptimeString := fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	return uptimeString
}
