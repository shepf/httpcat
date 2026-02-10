package midware

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"httpcat/internal/common"
	"httpcat/internal/common/utils"
	"httpcat/internal/common/ylog"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/rs/xid"
)

var whiteUrlList = []string{
	"/api/v1/file/upload",        //上传文件，我们需要开放出来，使用ak、sk方式生成专门的上传token，不支持界面操作上传
	"/api/v1/file/download",      // 下载文件，我们需要开放出来
	"/api/v1/imageManage/upload", // 上传图片，使用 UploadToken 校验
	"/api/v1/imageManage/download", // 下载图片
	"/api/v1/user/login/account"}

type AuthClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const (
	JWTExpireMinute = 60
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
		// 校验令牌到期时间
		expirationTime := claims["exp"].(float64)
		currentTime := time.Now().Unix()
		if currentTime > int64(expirationTime) {
			ylog.Errorf("VerifyToken", "Token expired")
			return nil, errors.New("Token expired")
		}

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

		//token := c.GetHeader("token")
		//if token == "" {
		//	ylog.Errorf("AuthRequired", "token is empty")
		//	c.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		// 因为前端：const authHeader = { Authorization: Bearer ${token} }，所以：
		token := c.GetHeader("Authorization")
		if token == "" {
			ylog.Errorf("AuthRequired", "token is empty")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 去掉 "Bearer " 前缀，只保留 token
		if len(token) > 7 && strings.ToUpper(token[0:7]) == "BEARER " {
			token = token[7:]
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

// tryAKSKAuth 尝试使用 AK/SK 签名认证，成功返回 true，未携带 AK/SK 头返回 false
func tryAKSKAuth(c *gin.Context) bool {
	akHeader := c.Request.Header.Get("AccessKey")
	sign := c.Request.Header.Get("Signature")
	timeStamp := c.Request.Header.Get("TimeStamp")

	// 没有 AK/SK 头，说明不是 Open API 请求
	if akHeader == "" || sign == "" || timeStamp == "" {
		return false
	}

	// 校验时间戳
	iTime, err := strconv.ParseInt(timeStamp, 10, 64)
	if err != nil {
		abort(c, fmt.Sprintf("TimeStamp Error: %s", err.Error()))
		return true
	}
	timeDiff := time.Now().Unix() - iTime
	if timeDiff >= 60 || timeDiff <= -60 {
		abort(c, fmt.Sprintf("timestamp out of range (±60s), diff=%ds", timeDiff))
		return true
	}

	// 校验签名
	sk := getSecKey(akHeader)
	if sk == "" {
		abort(c, "invalid AccessKey")
		return true
	}

	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		abort(c, err.Error())
		return true
	}
	_ = c.Request.Body.Close()
	c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	serverSign := generateSign(c.Request.Method, formatURLPath(c.Request.URL.Path), c.Request.URL.RawQuery, akHeader, timeStamp, sk, requestBody)
	if !signEqual(serverSign, sign) {
		ylog.Errorf("tryAKSKAuth", "signature mismatch for ak=%s method=%s path=%s", akHeader, c.Request.Method, c.Request.URL.Path)
		abort(c, "signature mismatch")
		return true
	}

	// AK/SK 认证成功，设置用户为 ak 标识
	c.Set("user", fmt.Sprintf("openapi:%s", akHeader))
	c.Next()
	return true
}

// TokenOrAKSKAuth 合并认证中间件：先尝试 AK/SK，不满足则走 JWT
func TokenOrAKSKAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 白名单豁免
		if utils.Contains(whiteUrlList, c.Request.URL.Path) {
			c.Next()
			return
		}

		// 如果开启了 Open API，优先尝试 AK/SK 认证
		if common.OpenAPIEnable {
			if tryAKSKAuth(c) {
				return // AK/SK 已处理（成功或失败都已响应）
			}
		}

		// 回退到 JWT Token 认证
		token := c.GetHeader("Authorization")
		if token == "" {
			ylog.Errorf("AuthRequired", "token is empty")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if len(token) > 7 && strings.ToUpper(token[0:7]) == "BEARER " {
			token = token[7:]
		}

		var userName string

		if strings.HasPrefix(token, "seesion-") {
			// session token (Redis, 当前已注释)
		} else {
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

		c.Set("user", userName)
		c.Next()
		return
	}
}
