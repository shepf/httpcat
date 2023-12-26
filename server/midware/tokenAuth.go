package midware

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"gin_web_demo/server/common"
	"gin_web_demo/server/common/utils"
	"gin_web_demo/server/common/ylog"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/rs/xid"
	"io"
	"net/http"
	"strings"
	"time"
)

var whiteUrlList = []string{
	"/api/v1/file/upload", //上传文件，我们需要开放出来，使用ak、sk方式生成专门的上传token，不支持界面操作上传
	"/api/v1/file/listFiles",
	"/api/v1/user/login/account"}

type AuthClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const (
	JWTExpireMinute = 720
)

var APITokenSecret = []byte(common.JwtSecret)

func CreateToken(payload jwt.Claims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string, secret []byte) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("parser failed")
		}

		return secret, nil
	})

	if err != nil {
		ylog.Errorf("VerifyToken", err.Error())
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}
	return nil, errors.New("verify token failed")
}

func checkPassword(password, salt, hash string) bool {
	t := sha1.New()
	_, err := io.WriteString(t, password+salt)
	if err != nil {
		return false
	}
	if fmt.Sprintf("%x", t.Sum(nil)) == hash {
		return true
	}
	return false
}

func GenPassword(password, salt string) string {
	t := sha1.New()
	_, err := io.WriteString(t, password+salt)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", t.Sum(nil))
}

func CheckUser(username, password string) (*common.User, error) {
	u := common.GetUser(username)
	if u == nil {
		ylog.Errorf("CheckUser", "user not found")
		return nil, errors.New("user not found")
	}

	if !checkPassword(password, u.Salt, u.Password) {
		return u, errors.New("verify password failed")
	}
	return u, nil
}

func GeneralJwtToken(userName string) (string, error) {
	return CreateToken(AuthClaims{
		Username: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(JWTExpireMinute * time.Minute).Unix(),
		},
	}, APITokenSecret)
}

func GeneralSession() string {
	return fmt.Sprintf("seesion-%s-%s", xid.New(), xid.New())
}

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//// infra.ApiAuth 变量为 false，即认证未开启，那么就不需要进行 Token 鉴权了。可
		//以直接调用 c.Next() 将控制权传递给下一个中间件或处理器，终止当前中间件的执行。
		//if !infra.ApiAuth {
		//	c.Next()
		//	return
		//}

		//url_whitelist
		if utils.Contains(whiteUrlList, c.Request.URL.Path) {
			c.Next()
			return
		}

		token := c.GetHeader("token")
		if token == "" {
			ylog.Errorf("AuthRequired", "token is empty")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var userName string

		if strings.HasPrefix(token, "seesion-") { // 如果 Token 前缀是 "seesion-"，则将其作为 key 从 Redis 中获取对应的用户名 userName。

			//userName = infra.Grds.Get(context.Background(), token).Val()
			//if userName == "" {
			//	c.AbortWithStatus(http.StatusUnauthorized)
			//	return
			//}
			//
			//err := infra.Grds.Expire(context.Background(), token, time.Duration(login.GetLoginSessionTimeoutMinute())*time.Minute).Err()
			//if err != nil {
			//	ylog.Errorf("TokenAuth", "Expire error %s", err.Error())
			//}
		} else { // 否则使用 JWT 验证 Token 的有效性，并从负载中获取用户名信息
			//jwt
			payload, err := VerifyToken(token, APITokenSecret)
			if err != nil {
				ylog.Errorf("AuthRequired", err.Error())
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			if payload == nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			currentUser, ok := (*payload)["username"]
			if currentUser == "" || !ok {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			userName = currentUser.(string)
		}

		// c.Set("user", userName)：将上下文中的键值对 "user" 设置为 userName。
		//这样，在后续的处理器函数或中间件中可以通过 c.Get("user") 方法获取到该值，用于后续的业务逻辑处理。
		c.Set("user", userName)
		c.Next()
		return
	}
}
