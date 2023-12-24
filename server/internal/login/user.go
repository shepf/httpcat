// Package login implements all login management interfaces.
package login

import (
	"database/sql"
	"fmt"
	"gin_web_demo/server/common"
	"gin_web_demo/server/common/ylog"
	"log"
	"sync"
	"time"
)

type User struct {
	Name        string     `json:"username" bson:"username"`
	Avatar      string     `json:"avatar" bson:"avatar"`
	UserID      string     `json:"userid" bson:"userid"`
	Email       string     `json:"email" bson:"email"`
	Signature   string     `json:"signature" bson:"signature"`
	Title       string     `json:"title" bson:"title"`
	Group       string     `json:"group" bson:"group"`
	Tags        []UserTag  `json:"tags" bson:"tags"`
	NotifyCount int        `json:"notify_count" bson:"notify_count"`
	UnreadCount int        `json:"unread_count" bson:"unread_count"`
	Country     string     `json:"country" bson:"country"`
	Access      []string   `json:"access" bson:"access"`
	Province    UserRegion `json:"province" bson:"province"`
	City        UserRegion `json:"city" bson:"city"`
	Address     string     `json:"address" bson:"address"`
	Phone       string     `json:"phone" bson:"phone"`

	Password           string       `json:"password" bson:"password"`
	PasswordUpdateTime int64        `json:"password_update_time" bson:"password_update_time"`
	Salt               string       `json:"salt" bson:"salt"`
	Level              int          `json:"level" bson:"level"`
	Config             []UserConfig `json:"config" bson:"config"`
}

type UserTag struct {
	Key   string `json:"key" bson:"key"`
	Label string `json:"label" bson:"label"`
}

type UserRegion struct {
	Label string `json:"label" bson:"label"`
	Key   string `json:"key" bson:"key"`
}

type UserConfig struct {
	Workspace string              `json:"workspace" bson:"workspace"`
	Favor     map[string][]string `json:"favor" bson:"favor"`
}

// 权限等级 0-->admin; 1-->高级用户(xxx)； 2-->xxx；
var (
	UserTable map[string]*User
	UserLock  sync.RWMutex
)

const LoginSessionTimeoutMin = 120

// GetUser find the user in the cache and returns, if user not exist, return nil.
// This interface is high-performance, but may not be up-to-date.
func GetUser(userName string) *User {
	UserLock.RLock()
	defer UserLock.RUnlock()
	user, ok := UserTable[userName]
	if !ok {
		return nil
	}
	return user
}

// GetLoginSessionTimeoutMinute returns the login session idle timeout time in minutes.
func GetLoginSessionTimeoutMinute() int64 {
	return LoginSessionTimeoutMin
}

// 用于初始化用户数据
// 在程序启动时从数据库加载用户数据，并将加载到的数据赋值给全局变量 UserTable，然后开启一个定时任务，在每隔3秒钟重新加载一次数据库中的用户数据，并更新全局变量 UserTable。
// 使用读写锁，它保证了在更新过程中不会发生竞争条件，从而确保了数据的一致性和安全性
func initUser() {
	table := loadUserFromDB()
	if table != nil {
		UserLock.Lock()
		UserTable = table
		UserLock.Unlock()
	}

	go func() {
		for {
			time.Sleep(3 * time.Second)
			table := loadUserFromDB()
			if table != nil {
				UserLock.Lock()
				UserTable = table
				UserLock.Unlock()
			}
		}
	}()
}

func loadUserFromDB() map[string]*User {
	if common.EnableSqlite {
		ylog.Infof("loadUserFromDB", "loadUserFromDB")
		// 读取 SQLite 数据库文件路径配置项
		dbPath := common.SqliteDBPath

		// 打开 SQLite 数据库连接
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			ylog.Errorf("uploadFile", "open db failed, err:%v", err)
			return nil
		}
		defer db.Close()

		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			log.Fatalf("loadUserFromDB: %v", err)
			return nil
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			log.Fatalf("Error getting column names: %v", err)
			return nil
		}

		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for rows.Next() {
			for i := range columns {
				valuePtrs[i] = &values[i]
			}
			if err := rows.Scan(valuePtrs...); err != nil {
				log.Fatalf("Error scanning row: %v", err)
				return nil
			}
			for i, col := range columns {
				val := values[i]
				switch val.(type) {
				case nil:
					fmt.Println(col, ": NULL")
				case []byte:
					fmt.Println(col, ": ", string(val.([]byte)))
				default:
					fmt.Println(col, ": ", val)
				}
			}
		}
		if err := rows.Err(); err != nil {
			log.Fatalf("Error in rows iteration: %v", err)
			return nil
		}

		userTable := map[string]*User{}
		for rows.Next() {
			var user User
			err := rows.Scan(&user.Name, &user.Email) // 假设用户表中有用户名和电子邮件字段
			if err != nil {
				log.Printf("loadUserFromDB: %v", err)
				continue
			}
			userTable[user.Name] = &user
		}
		return userTable

	}

	return nil
}
