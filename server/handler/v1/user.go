package v1

import (
	"gin_web_demo/server/common"
	"gin_web_demo/server/common/ylog"
	"gin_web_demo/server/midware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

type AuthRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Type      string `json:"type" binding:"required"`
	AutoLogin bool   `json:"autoLogin"`
}

func UserLogin(c *gin.Context) {
	var user AuthRequest

	// 通过c.BindJSON(&user)将请求体中的JSON数据绑定到user结构体中。如果绑定失败，会返回参数无效的错误响应。
	err := c.BindJSON(&user)
	if err != nil {
		common.CreateResponse(c, common.ParamInvalidErrorCode, err.Error())
		return
	}

	if user.Username == "admin" {
		//使用jwt token
		_, err := midware.CheckUser(user.Username, user.Password)
		if err != nil {
			common.CreateResponse(c, common.AuthFailedErrorCode, err.Error())
			return
		}
		token, err := midware.GeneralJwtToken(user.Username)
		if err != nil {
			common.CreateResponse(c, common.AuthFailedErrorCode, err.Error())
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
		common.CreateResponse(c, common.AuthFailedErrorCode, "verify password failed")
	} else {
		token := midware.GeneralSession()
		//todo token存储到reids
		//err = infra.Grds.Set(context.Background(), token, user.Username, time.Duration(login.GetLoginSessionTimeoutMinute())*time.Minute).Err()
		if err != nil {
			ylog.Errorf("UserLogin", "Set %s redis error %s", user.Username, err.Error())
		}
		common.CreateResponse(c, common.SuccessCode, bson.M{"token": token})
	}
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
