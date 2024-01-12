package v1

import (
	"fmt"
	"gin_web_demo/server/common"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func GetDownloadStatistics(c *gin.Context) {

	common.CreateResponse(c, common.SuccessCode, nil)

}

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

	var monthPercentage string
	if lastMonthUploadCount == 0 {
		monthPercentage = "0%"
	} else {
		percentage := float64(monthUploadCount-lastMonthUploadCount) / float64(lastMonthUploadCount) * 100
		monthPercentage = formatPercentage(percentage)
	}

	// 将统计信息返回给前端
	common.CreateResponse(c, common.SuccessCode, gin.H{
		"todayUploadCount": todayUploadCount,
		"todayPercentage":  todayPercentage,
		"monthUploadCount": monthUploadCount,
		"monthPercentage":  monthPercentage,
		"totalUploadCount": totalUploadCount,
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
	}
	return fmt.Sprintf("%s%.2f%%", sign, percentage)
}
