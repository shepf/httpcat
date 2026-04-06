package common

import (
	"httpcat/internal/common/ylog"
	"httpcat/internal/models"
	"time"

	"gorm.io/gorm"
)

func InitializeShareTable(db *gorm.DB) {
	err := db.AutoMigrate(&models.ShareModel{})
	if err != nil {
		ylog.Errorf("initDB", "create t_share table failed, err:%v", err)
		return
	}

	// 启动过期分享定时清理
	go startShareCleanup()
}

// startShareCleanup 定时清理过期分享（每小时执行一次）
func startShareCleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// 启动时先执行一次
	cleanupExpiredShares()

	for range ticker.C {
		cleanupExpiredShares()
	}
}

// cleanupExpiredShares 将过期或达到下载上限的分享标记为非活跃
func cleanupExpiredShares() {
	db, err := GetDB()
	if err != nil {
		return
	}

	now := time.Now()

	// 将已过期的分享设为 is_active=false
	result := db.Model(&models.ShareModel{}).
		Where("is_active = ? AND expire_at IS NOT NULL AND expire_at < ?", true, now).
		Update("is_active", false)
	if result.RowsAffected > 0 {
		ylog.Infof("ShareCleanup", "cleaned up %d expired shares", result.RowsAffected)
	}

	// 将达到下载上限的分享设为 is_active=false
	result = db.Model(&models.ShareModel{}).
		Where("is_active = ? AND max_downloads > 0 AND cur_downloads >= max_downloads", true).
		Update("is_active", false)
	if result.RowsAffected > 0 {
		ylog.Infof("ShareCleanup", "cleaned up %d max-download-reached shares", result.RowsAffected)
	}
}
