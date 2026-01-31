package common

import (
	"fmt"
	"gin_web_demo/server/common/ylog"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	dbOnce     sync.Once
)

// GetDB 获取数据库连接实例（单例模式）
func GetDB() (*gorm.DB, error) {
	if !EnableSqlite {
		return nil, fmt.Errorf("SQLite is not enabled")
	}

	var err error
	dbOnce.Do(func() {
		dbInstance, err = gorm.Open(sqlite.Open(SqliteDBPath), &gorm.Config{})
		if err != nil {
			ylog.Errorf("GetDB", "Failed to open database: %v", err)
		}
	})

	if err != nil {
		return nil, err
	}

	return dbInstance, nil
}

// ResetDB 重置数据库连接（用于测试）
func ResetDB() {
	dbOnce = sync.Once{}
	dbInstance = nil
}
