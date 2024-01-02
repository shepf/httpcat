package common

import (
	"crypto/sha1"
	"fmt"
	"gin_web_demo/server/common/ylog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"log"
	"sync"
	"time"
)

type User struct {
	ID                 uint   `gorm:"primary_key"`
	Username           string `gorm:"column:username;unique"`
	Avatar             string `gorm:"column:avatar"`
	UserID             string `gorm:"column:userid"`
	Email              string `gorm:"column:email"`
	Signature          string `gorm:"column:signature"`
	Title              string `gorm:"column:title"`
	Group              string `gorm:"column:group"`
	Tags               []byte `gorm:"column:tags"`
	NotifyCount        int    `gorm:"column:notify_count"`
	UnreadCount        int    `gorm:"column:unread_count"`
	Country            string `gorm:"column:country"`
	Access             []byte `gorm:"column:access"`
	Province           []byte `gorm:"column:province"`
	City               []byte `gorm:"column:city"`
	Address            string `gorm:"column:address"`
	Phone              string `gorm:"column:phone"`
	Password           string `gorm:"column:password"`
	PasswordUpdateTime int64  `gorm:"column:password_update_time"`
	Salt               string `gorm:"column:salt"`
	Level              int    `gorm:"column:level"`
	Config             []byte `gorm:"column:config"`
}

type UserTag struct {
	Key   string `json:"key" bson:"key"`
	Label string `json:"label" bson:"label"`
}

type UserRegion struct {
	Label string `json:"label" bson:"label"`
	Key   string `json:"key" bson:"key"`
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
func InitUser() {
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
	if EnableSqlite {
		ylog.Infof("loadUserFromDB", "loadUserFromDB")
		// 读取 SQLite 数据库文件路径配置项
		dbPath := SqliteDBPath

		// 打开 SQLite 数据库连接
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			ylog.Errorf("initDB", "open db failed, err:%v", err)
		}

		var users []User
		if err := db.Find(&users).Error; err != nil {
			log.Fatalf("loadUserFromDB: %v", err)
			return nil
		}
		userTable := make(map[string]*User)
		for _, user := range users {
			userTable[user.Username] = &user
		}

		return userTable

	}

	return nil
}

func InitializeUserTable(db *gorm.DB) {

	err := db.AutoMigrate(&User{})
	if err != nil {
		ylog.Errorf("initDB", "create users table failed, err:%v", err)
		return
	}

	var count int64
	db.Model(&User{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		t := sha1.New()
		password := "admin"
		salt := "sss"
		_, err := io.WriteString(t, password+salt)
		hashString := fmt.Sprintf("%x", t.Sum(nil))

		// 插入默认记录
		user := User{
			Username:           "admin",
			Avatar:             "https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png",
			UserID:             "00000001",
			Email:              "antdesign@alipay.com",
			Signature:          "海纳百川，有容乃大",
			Title:              "交互专家",
			Group:              "蚂蚁金服－某某某事业群－某某平台部－某某技术部－UED",
			Tags:               []byte("[{\"key\":\"0\",\"label\":\"很有想法的\"},{\"key\":\"1\",\"label\":\"专注设计\"},{\"key\":\"2\",\"label\":\"辣~\"},{\"key\":\"3\",\"label\":\"大长腿\"},{\"key\":\"4\",\"label\":\"川妹子\"},{\"key\":\"5\",\"label\":\"海纳百川\"}]"),
			NotifyCount:        12,
			UnreadCount:        11,
			Country:            "China",
			Access:             []byte("ss"),
			Province:           []byte("{\"label\":\"浙江省\",\"key\":\"330000\"}"),
			City:               []byte("{\"label\":\"杭州市\",\"key\":\"330100\"}"),
			Address:            "西湖区工专路 77 号",
			Phone:              "0752-268888888",
			Password:           hashString, // 根据具体逻辑设置密码
			PasswordUpdateTime: 0,
			Salt:               salt, // 根据具体逻辑设置盐值
			Level:              0,
			Config:             nil, // 根据具体逻辑设置配置信息
		}
		err = db.Create(&user).Error
		if err != nil {
			ylog.Errorf("initDB", "insert user record failed, err:%v", err)
			return
		}
	}
}

type UploadTokenItem struct {
	ID        int       `gorm:"primary_key" json:"id"`
	Appkey    string    `gorm:"column:appkey;unique" json:"appkey"`
	Appsecret string    `gorm:"column:app_secret" json:"appsecret"`
	State     string    `gorm:"column:state" json:"state"`
	Desc      string    `gorm:"column:desc" json:"desc"`
	CreatedAt time.Time `gorm:"column:create_at" json:"created_at"`
}

// 指定表名
func (UploadTokenItem) TableName() string {
	return "t_upload_token"
}

func InitializeUploadTokenTable(db *gorm.DB) {
	err := db.AutoMigrate(&UploadTokenItem{})
	if err != nil {
		ylog.Errorf("initDB", "create upload_tokens table failed, err:%v", err)
		return
	}

	var count int64
	db.Model(&UploadTokenItem{}).Where("appkey = ?", "httpcat").Count(&count)
	if count == 0 {
		// 插入默认记录
		token := UploadTokenItem{
			Appkey:    AppKey,
			Appsecret: AppSecret,
			State:     "open",
			Desc:      "系统初始化默认appkey:httpcat",
			CreatedAt: time.Now(),
		}
		err = db.Create(&token).Error
		if err != nil {
			ylog.Errorf("initDB", "insert upload token record failed, err:%v", err)
			return
		}
	}
}
