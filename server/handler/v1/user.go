package v1

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"gin_web_demo/server/common"
	"gin_web_demo/server/common/ylog"
	"gin_web_demo/server/midware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type AuthRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Type      string `json:"type" binding:"required"`
	AutoLogin bool   `json:"autoLogin"`
}

func UserLogin(c *gin.Context) {
	ylog.Infof("UserLogin", "UserLogin function called")

	var user AuthRequest

	// 通过c.BindJSON(&user)将请求体中的JSON数据绑定到user结构体中。如果绑定失败，会返回参数无效的错误响应。
	err := c.BindJSON(&user)
	if err != nil {
		common.BadRequest(c, err.Error())
		return
	}

	if user.Username == "admin" {
		//使用jwt token
		_, err := midware.CheckUser(user.Username, user.Password)
		if err != nil {
			common.Unauthorized(c, err.Error())
			return
		}
		token, err := midware.GeneralJwtToken(user.Username)
		if err != nil {
			common.Unauthorized(c, err.Error())
			return
		}

		c.JSON(
			http.StatusOK,
			bson.M{"token": token,
				"currentAuthority": "access", "type": "account", "status": "ok"},
		)
		return
	}

	_, err = midware.CheckUser(user.Username, user.Password)
	//密码校验
	if err != nil {
		common.Unauthorized(c, "verify password failed")
	} else {
		token := midware.GeneralSession()
		//todo token存储到reids
		//err = infra.Grds.Set(context.Background(), token, user.Username, time.Duration(login.GetLoginSessionTimeoutMinute())*time.Minute).Err()
		if err != nil {
			ylog.Errorf("UserLogin", "Set %s redis error %s", user.Username, err.Error())
		}
		common.CreateResponse(c, common.SuccessCode, bson.M{"token": token})
	}

	ylog.Infof("UserLogin", "UserLogin function completed")

}

func ChangePasswd(c *gin.Context) {

	var params struct {
		OldPassword string `form:"oldPassword" binding:"required"`
		NewPassword string `form:"newPassword" binding:"required"`
	}

	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前登录用户名, 当前只返回admin
	username, ok := c.Get("user")
	if !ok {
		common.CreateResponse(c, common.ErrorCode, "Failed to get user information")
		return
	}

	// 从数据库中查询用户信息
	var user common.User
	dbPath := common.SqliteDBPath
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		// 处理数据库连接错误
		common.CreateResponse(c, common.ErrorCode, err)
		return
	}

	db.Where("username = ?", username.(string)).First(&user)

	// 验证旧密码是否正确
	t := sha1.New()
	_, err = io.WriteString(t, params.OldPassword+user.Salt)
	oldPasswordHash := fmt.Sprintf("%x", t.Sum(nil))
	if err != nil || oldPasswordHash != user.Password {
		common.Unauthorized(c, "Incorrect old password")
		return
	}

	// 根据具体的逻辑，生成新密码的哈希值和盐值
	t = sha1.New()
	_, err = io.WriteString(t, params.NewPassword+user.Salt)
	newPasswordHash := fmt.Sprintf("%x", t.Sum(nil))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new password hash"})
		return
	}

	// 更新用户的密码信息
	user.Password = newPasswordHash
	user.PasswordUpdateTime = time.Now().Unix()

	// 保存更新后的用户信息到数据库
	db.Save(&user)

	// 返回成功响应
	common.CreateResponse(c, common.SuccessCode, "Password changed successfully")

}

type UserInfoVO struct {
	ID          uint   `json:"id"`
	Username    string `json:"name"`
	Avatar      string `json:"avatar"`
	UserID      string `json:"userid"`
	Email       string `json:"email"`
	Signature   string `json:"signature"`
	Title       string `json:"title"`
	Group       string `json:"group"`
	NotifyCount int    `json:"notifyCount"`
	UnreadCount int    `json:"unreadCount"`
	Country     string `json:"country"`
	Address     string `json:"address"`
	Access      string `json:"access"`
	Phone       string `json:"phone"`

	Level int `json:"level"`
}

func UserInfo(c *gin.Context) {

	// 获取当前登录用户名, 当前只返回admin
	username, ok := c.Get("user")
	if !ok {
		common.CreateResponse(c, common.ErrorCode, "Failed to get user information")
		return
	}

	user := common.GetUser(username.(string))

	// 构建包含需要保留字段的新结构体
	info := UserInfoVO{
		ID:          user.ID,
		Username:    user.Username,
		Avatar:      user.Avatar,
		UserID:      user.UserID,
		Email:       user.Email,
		Signature:   user.Signature,
		Title:       user.Title,
		Group:       user.Group,
		NotifyCount: user.NotifyCount,
		UnreadCount: user.UnreadCount,
		Country:     user.Country,
		Address:     user.Address,
		Access:      "admin",
		Phone:       user.Phone,
		Level:       user.Level,
	}

	common.CreateResponse(c, common.SuccessCode, info)
}

func UserLoginout(c *gin.Context) {
	token := c.GetHeader("token")
	ylog.Infof("UserLoginout", "delete token: %s", token)

	// TODO token从redis中删除
	//err := infra.Grds.Del(context.Background(), token).Err()

	//if err != nil {
	//	common.CreateResponse(c, common.SuccessCode, err.Error())
	//} else {
	//	common.CreateResponse(c, common.SuccessCode, nil)
	//}

	common.CreateResponse(c, common.SuccessCode, nil)
}

func UploadAvatar(c *gin.Context) {
	// 获取当前登录用户名
	username, _ := c.Get("user")
	uploadedFile, err := c.FormFile("avatar")
	if err != nil {
		// 处理文件上传错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 打开上传的文件
	file, err := uploadedFile.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// 读取文件内容
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 将文件内容进行 Base64 编码
	encodedString := base64.StdEncoding.EncodeToString(fileBytes)
	// 更新用户的 Avatar 属性
	dbPath := common.SqliteDBPath
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		// 处理数据库连接错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var user common.User
	db.Where("username = ?", username.(string)).First(&user)
	user.Avatar = encodedString
	db.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Avatar uploaded successfully"})

}
