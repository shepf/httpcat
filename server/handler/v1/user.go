package v1

import (
	"gin_web_demo/server/common"
	"gin_web_demo/server/common/ylog"
	"gin_web_demo/server/midware"
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
		common.CreateResponse(c, common.SuccessCode, map[string]interface{}{"token": token})
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
