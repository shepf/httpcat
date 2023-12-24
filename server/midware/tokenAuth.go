package midware

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"gin_web_demo/server/common"
	"gin_web_demo/server/common/ylog"
	"gin_web_demo/server/internal/login"
	"github.com/golang-jwt/jwt"
	"github.com/rs/xid"
	"io"
	"time"
)

var whiteUrlList = []string{
	"/api/v1/agent/heartbeat/evict",
	"/api/v1/user/login"}

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

func CheckUser(username, password string) (*login.User, error) {
	u := login.GetUser(username)
	if u == nil {
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
