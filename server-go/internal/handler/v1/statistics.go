package v1

import (
	"fmt"
	"httpcat/internal/common"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"math"
	"net/http"
	"time"
)

func GetUploadStatistics(c *gin.Context) {
	dbPath := common.SqliteDBPath
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	db.Debug()

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
	dbPath := common.SqliteDBPath
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	db.Debug()

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
