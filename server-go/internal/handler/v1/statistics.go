package v1

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"httpcat/internal/common"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUploadStatistics(c *gin.Context) {
	db, err := common.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 统计信息
	todayUploadCount := getTodayUploadCount(db)
	yesterdayUploadCount := getYesterdayUploadCount(db)
	monthUploadCount := getMonthUploadCount(db)
	lastMonthUploadCount := getLastMonthUploadCount(db)
	totalUploadCount := getTotalUploadCount(db)

	// 计算百分比
	var todayPercentage string
	if yesterdayUploadCount == 0 {
		todayPercentage = "0%"
	} else {
		percentage := float64(todayUploadCount-yesterdayUploadCount) / float64(yesterdayUploadCount) * 100
		todayPercentage = formatPercentage(percentage)
	}
	fmt.Println("todayPercentage:", todayPercentage)

	var monthPercentage string
	if lastMonthUploadCount == 0 {
		monthPercentage = "0%"
	} else {
		percentage := float64(monthUploadCount-lastMonthUploadCount) / float64(lastMonthUploadCount) * 100
		monthPercentage = formatPercentage(percentage)
	}

	// 将统计信息返回给前端
	common.CreateResponse(c, common.SuccessCode, gin.H{
		"todayUploadCount":     todayUploadCount,
		"yesterdayUploadCount": yesterdayUploadCount,
		"todayPercentage":      todayPercentage,
		"monthUploadCount":     monthUploadCount,
		"lastMonthUploadCount": lastMonthUploadCount,
		"monthPercentage":      monthPercentage,
		"totalUploadCount":     totalUploadCount,
	})

}

func getTodayUploadCount(db *gorm.DB) int64 {
	var count int64
	today := time.Now().Format("2006-01-02")
	db.Table("t_upload_log").Where("DATE(upload_time) = ?", today).Count(&count)
	return count
}

func getYesterdayUploadCount(db *gorm.DB) int64 {
	var count int64
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	db.Table("t_upload_log").Where("DATE(upload_time) = ?", yesterday).Count(&count)
	return count
}

func getMonthUploadCount(db *gorm.DB) int64 {
	var count int64
	db.Table("t_upload_log").Where("strftime('%Y-%m', upload_time) = strftime('%Y-%m', 'now')").Count(&count)
	return count
}

func getLastMonthUploadCount(db *gorm.DB) int64 {
	var count int64
	db.Table("t_upload_log").Where("strftime('%Y-%m', upload_time) = strftime('%Y-%m', 'now', '-1 month')").Count(&count)
	return count
}

func getTotalUploadCount(db *gorm.DB) int64 {
	var count int64
	db.Table("t_upload_log").Count(&count)
	return count
}

func formatPercentage(percentage float64) string {
	sign := ""
	if percentage > 0 {
		sign = "+"
	} else if percentage < 0 {
		sign = "-"
	} else {
		return "0%"
	}
	// 去掉 percentage 本身的负数，格式化时候再加上，否则会出现 --xx% 的情况
	return fmt.Sprintf("%s%.2f%%", sign, math.Abs(percentage))
}

// Path: server\handler\v1\statistics.go
func GetDownloadStatistics(c *gin.Context) {
	db, err := common.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 统计信息
	todayDownloadCount := getTodayDownloadCount(db)
	yesterdayDownloadCount := getYesterdayDownloadCount(db)
	monthDownloadCount := getMonthDownloadCount(db)
	lastMonthDownloadCount := getLastMonthDownloadCount(db)
	totalDownloadCount := getTotalDownloadCount(db)

	// 计算百分比
	var todayPercentage string
	if yesterdayDownloadCount == 0 {
		todayPercentage = "0%"
	} else {
		percentage := float64(todayDownloadCount-yesterdayDownloadCount) / float64(yesterdayDownloadCount) * 100
		todayPercentage = formatPercentage(percentage)
	}

	var monthPercentage string
	if lastMonthDownloadCount == 0 {
		monthPercentage = "0%"
	} else {
		percentage := float64(monthDownloadCount-lastMonthDownloadCount) / float64(lastMonthDownloadCount) * 100
		monthPercentage = formatPercentage(percentage)
	}

	// 将统计信息返回给前端
	common.CreateResponse(c, common.SuccessCode, gin.H{
		"todayDownloadCount":     todayDownloadCount,
		"yesterdayDownloadCount": yesterdayDownloadCount, // 昨天下载量
		"todayPercentage":        todayPercentage,
		"monthDownloadCount":     monthDownloadCount,
		"lastMonthDownloadCount": lastMonthDownloadCount, // 上个月下载量
		"monthPercentage":        monthPercentage,
		"totalDownloadCount":     totalDownloadCount,
	})
}

func getTodayDownloadCount(db *gorm.DB) int64 {
	var count int64
	today := time.Now().Format("2006-01-02")
	db.Model(&common.DownloadLogModel{}).Where("download_time >= ?", today).Count(&count)
	return count
}

func getYesterdayDownloadCount(db *gorm.DB) int64 {
	var count int64
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	today := time.Now().Format("2006-01-02")
	db.Model(&common.DownloadLogModel{}).Where("download_time >= ? AND download_time < ?", yesterday, today).Count(&count)
	return count
}

func getMonthDownloadCount(db *gorm.DB) int64 {
	var count int64
	firstDayOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1).Format("2006-01-02")
	db.Model(&common.DownloadLogModel{}).Where("download_time >= ?", firstDayOfMonth).Count(&count)
	return count
}

func getLastMonthDownloadCount(db *gorm.DB) int64 {
	var count int64
	firstDayOfLastMonth := time.Now().AddDate(0, -1, -time.Now().Day()+1).Format("2006-01-02")
	lastDayOfLastMonth := time.Now().AddDate(0, 0, -time.Now().Day()).Format("2006-01-02")
	db.Model(&common.DownloadLogModel{}).Where("download_time >= ? AND download_time <= ?", firstDayOfLastMonth, lastDayOfLastMonth).Count(&count)
	return count
}

func getTotalDownloadCount(db *gorm.DB) int64 {
	var count int64
	db.Model(&common.DownloadLogModel{}).Count(&count)
	return count
}

// GetDownloadHistoryLogs 获取下载历史记录（分页）
func GetDownloadHistoryLogs(c *gin.Context) {
	var params struct {
		Current  int    `form:"current" binding:"required"`
		PageSize int    `form:"pageSize" binding:"required"`
		FileName string `form:"filename"`
		FileMD5  string `form:"file_md5"`
		IP       string `form:"ip"`
	}
	if err := c.ShouldBindQuery(&params); err != nil {
		common.CreateResponse(c, common.ParamInvalidErrorCode, err.Error())
		return
	}

	db, err := common.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	offset := (params.Current - 1) * params.PageSize
	query := db.Model(&common.DownloadLogModel{}).Offset(offset).Limit(params.PageSize).Order("download_time DESC")
	if params.FileName != "" {
		query = query.Where("filename LIKE ?", "%"+params.FileName+"%")
	}
	if params.FileMD5 != "" {
		query = query.Where("file_md5 = ?", params.FileMD5)
	}
	if params.IP != "" {
		query = query.Where("ip = ?", params.IP)
	}

	var logs []common.DownloadLogModel
	if err := query.Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 查询总数
	var total int64
	countQuery := db.Model(&common.DownloadLogModel{})
	if params.FileName != "" {
		countQuery = countQuery.Where("filename LIKE ?", "%"+params.FileName+"%")
	}
	if params.FileMD5 != "" {
		countQuery = countQuery.Where("file_md5 = ?", params.FileMD5)
	}
	if params.IP != "" {
		countQuery = countQuery.Where("ip = ?", params.IP)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	common.CreateResponse(c, common.SuccessCode, gin.H{
		"list":     logs,
		"current":  params.Current,
		"pageSize": params.PageSize,
		"total":    total,
	})
}

// GetFileOverview 获取文件总览统计：文件总数、目录数、总大小
func GetFileOverview(c *gin.Context) {
	baseDir := common.GetDownloadDir()

	var totalFiles int64
	var totalDirs int64
	var totalSize int64

	// 递归统计
	filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 跳过无法访问的文件
		}
		if info.IsDir() {
			totalDirs++
		} else {
			totalFiles++
			totalSize += info.Size()
		}
		return nil
	})

	// 减去根目录本身
	if totalDirs > 0 {
		totalDirs--
	}

	common.CreateResponse(c, common.SuccessCode, gin.H{
		"totalFiles": totalFiles,
		"totalDirs":  totalDirs,
		"totalSize":  totalSize,
		"totalSizeFormatted": formatBytes(totalSize),
	})
}

func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)
	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	}
	return fmt.Sprintf("%d B", bytes)
}
