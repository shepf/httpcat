package v1

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"httpcat/internal/common"
	"httpcat/internal/common/ylog"
	"httpcat/internal/midware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type AuthRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Type      string `json:"type" binding:"required"`
	AutoLogin bool   `json:"autoLogin"`
}

func UserLogin(c *gin.Context) {
	ylog.Infof("UserLogin", "UserLogin function called")

	clientIP := c.ClientIP()

	var user AuthRequest

	// 通过c.BindJSON(&user)将请求体中的JSON数据绑定到user结构体中。如果绑定失败，会返回参数无效的错误响应。
	err := c.BindJSON(&user)
	if err != nil {
		common.BadRequest(c, err.Error())
		return
	}

	if user.Username == common.DefaultAdminUsername {
		//使用jwt token
		userInfo, err := midware.CheckUser(user.Username, user.Password)
		if err != nil {
			midware.RecordLoginFailure(clientIP)
			common.Unauthorized(c, err.Error())
			return
		}
		token, err := midware.GeneralJwtToken(user.Username)
		if err != nil {
			common.Unauthorized(c, err.Error())
			return
		}

		midware.RecordLoginSuccess(clientIP)

		c.JSON(
			http.StatusOK,
			bson.M{"token": token,
				"currentAuthority": "access", "type": "account", "status": "ok", "mustChangePassword": common.MustChangePassword(userInfo)},
		)
		return
	}

	userInfo, err := midware.CheckUser(user.Username, user.Password)
	//密码校验
	if err != nil {
		midware.RecordLoginFailure(clientIP)
		common.Unauthorized(c, "verify password failed")
	} else {
		midware.RecordLoginSuccess(clientIP)
		token := midware.GeneralSession()
		//todo token存储到reids
		//err = infra.Grds.Set(context.Background(), token, user.Username, time.Duration(login.GetLoginSessionTimeoutMinute())*time.Minute).Err()
		common.CreateResponse(c, common.SuccessCode, bson.M{"token": token, "mustChangePassword": common.MustChangePassword(userInfo)})
	}

	ylog.Infof("UserLogin", "UserLogin function completed")

}

func ChangePasswd(c *gin.Context) {

	var params struct {
		OldPassword string `form:"oldPassword" json:"oldPassword" binding:"required"`
		NewPassword string `form:"newPassword" json:"newPassword" binding:"required"`
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
	db, err := common.GetDB()
	if err != nil {
		// 处理数据库连接错误
		common.CreateResponse(c, common.ErrorCode, err)
		return
	}

	db.Where("username = ?", username.(string)).First(&user)

	valid, _, err := common.VerifyPassword(&user, params.OldPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify old password"})
		return
	}
	if !valid {
		common.Unauthorized(c, "Incorrect old password")
		return
	}

	newPassword := strings.TrimSpace(params.NewPassword)
	if newPassword == "" {
		common.BadRequest(c, "新密码不能为空")
		return
	}
	newPasswordHash, err := common.HashPassword(newPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new password hash"})
		return
	}

	// 更新用户的密码信息
	user.Password = newPasswordHash
	user.Salt = ""
	user.PasswordUpdateTime = time.Now().Unix()

	// 保存更新后的用户信息到数据库
	if err := db.Model(&user).Updates(map[string]interface{}{
		"password":             user.Password,
		"salt":                 user.Salt,
		"password_update_time": user.PasswordUpdateTime,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new password"})
		return
	}
	common.UpdateUserPasswordCache(user.Username, user.Password, user.Salt, user.PasswordUpdateTime)

	// 返回成功响应
	common.CreateResponse(c, common.SuccessCode, "Password changed successfully")

}

type UserInfoVO struct {
	ID                 uint   `json:"id"`
	Username           string `json:"name"`
	Avatar             string `json:"avatar"`
	UserID             string `json:"userid"`
	Email              string `json:"email"`
	Signature          string `json:"signature"`
	Title              string `json:"title"`
	Group              string `json:"group"`
	NotifyCount        int    `json:"notifyCount"`
	UnreadCount        int    `json:"unreadCount"`
	Country            string `json:"country"`
	Address            string `json:"address"`
	Access             string `json:"access"`
	Phone              string `json:"phone"`
	MustChangePassword bool   `json:"mustChangePassword"`

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
		ID:                 user.ID,
		Username:           user.Username,
		Avatar:             user.Avatar,
		UserID:             user.UserID,
		Email:              user.Email,
		Signature:          user.Signature,
		Title:              user.Title,
		Group:              user.Group,
		NotifyCount:        user.NotifyCount,
		UnreadCount:        user.UnreadCount,
		Country:            user.Country,
		Address:            user.Address,
		Access:             "admin",
		Phone:              user.Phone,
		MustChangePassword: common.MustChangePassword(user),
		Level:              user.Level,
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
	db, err := common.GetDB()
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
