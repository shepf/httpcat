package common

import (
	"os"
	"path/filepath"
	"time"

	"httpcat/internal/common/ylog"
	"httpcat/internal/models"

	"gorm.io/gorm"
)

// ChunkTempDir 返回分片临时存储根目录（位于 SQLite 同级 data/chunks）
// 例如 SqliteDBPath = "./data/httpcat_sqlite.db" -> "./data/chunks"
func ChunkTempDir() string {
	base := filepath.Dir(SqliteDBPath)
	if base == "" || base == "." {
		base = "./data"
	}
	return filepath.Join(base, "chunks")
}

// ChunkSessionDir 返回某个 uploadID 的分片临时目录
func ChunkSessionDir(uploadID string) string {
	return filepath.Join(ChunkTempDir(), uploadID)
}

// InitializeUploadSessionTable 初始化分片上传会话表（v0.7.0）
func InitializeUploadSessionTable(db *gorm.DB) {
	if err := db.AutoMigrate(&models.UploadSessionModel{}); err != nil {
		ylog.Errorf("initDB", "create t_upload_session table failed, err:%v", err)
		return
	}

	// 确保分片临时目录存在
	if err := os.MkdirAll(ChunkTempDir(), 0o755); err != nil {
		ylog.Errorf("initDB", "create chunks dir failed: %v", err)
	}

	// 启动过期会话清理任务
	go startUploadSessionCleanup()
}

// startUploadSessionCleanup 每 30 分钟清理一次过期/已完成但遗留分片的会话
func startUploadSessionCleanup() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	// 启动时先执行一次
	cleanupExpiredUploadSessions()

	for range ticker.C {
		cleanupExpiredUploadSessions()
	}
}

// cleanupExpiredUploadSessions 删除过期会话的分片目录并标记状态
func cleanupExpiredUploadSessions() {
	db, err := GetDB()
	if err != nil {
		return
	}

	now := time.Now()

	// 找出所有已过期但仍为 active 的会话
	var expiredSessions []models.UploadSessionModel
	if err := db.Where("status = ? AND expire_at < ?", "active", now).Find(&expiredSessions).Error; err != nil {
		ylog.Errorf("UploadSessionCleanup", "query expired sessions failed: %v", err)
		return
	}

	for _, s := range expiredSessions {
		dir := ChunkSessionDir(s.UploadID)
		if err := os.RemoveAll(dir); err != nil {
			ylog.Errorf("UploadSessionCleanup", "remove chunks dir %s failed: %v", dir, err)
			continue
		}
		db.Model(&models.UploadSessionModel{}).
			Where("upload_id = ?", s.UploadID).
			Update("status", "aborted")
	}

	if len(expiredSessions) > 0 {
		ylog.Infof("UploadSessionCleanup", "cleaned up %d expired upload sessions", len(expiredSessions))
	}

	// 清理已完成超过 1 天但遗留分片目录的会话（防御性清理）
	var staleCompleted []models.UploadSessionModel
	cutoff := now.Add(-24 * time.Hour)
	if err := db.Where("status = ? AND updated_at < ?", "completed", cutoff).Find(&staleCompleted).Error; err == nil {
		for _, s := range staleCompleted {
			dir := ChunkSessionDir(s.UploadID)
			if _, err := os.Stat(dir); err == nil {
				_ = os.RemoveAll(dir)
			}
		}
	}
}
