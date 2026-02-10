package midware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"httpcat/internal/common"
	"httpcat/internal/common/ylog"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/levigross/grequests"
)

var (
	ak = common.SvrAK
	sk = common.SvrSK
)

// getSecKey 根据 AccessKey 查找对应的 SecretKey
func getSecKey(ak string) string {
	sk, ok := common.HttpAkSkMap[ak]
	if !ok {
		return ""
	}
	return sk
}

// sha256Hex 计算字节数组的 SHA256 哈希值，返回十六进制字符串。
// 即使输入为空（nil 或长度为 0），也会计算空字符串的 SHA256 值。
// 这确保签名字符串末尾始终有确定值，避免客户端集成时的换行符陷阱（如 shell $() 会吞末尾换行）。
// 空 body 的 SHA256: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
func sha256Hex(in []byte) string {
	h := sha256.New()
	if len(in) > 0 {
		h.Write(in)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// generateSign 生成 HMAC-SHA256 签名。
// 签名字符串格式：{Method}\n{Path}\n{Query}\n{AK}\n{Timestamp}\n{BodySHA256}
// 各字段之间使用真正的换行符（0x0a）分隔，与 AWS Signature V4 等行业标准一致。
func generateSign(method, url, query, ak, timestamp, sk string, requestBody []byte) string {
	signStr := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", method, url, query, ak, timestamp, sha256Hex(requestBody))
	return hmacSha256(signStr, sk)
}

// hmacSha256 使用 HMAC-SHA256 算法对数据签名，返回十六进制字符串。
func hmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// signEqual 使用恒定时间比较两个签名，防止时序攻击。
func signEqual(a, b string) bool {
	return hmac.Equal([]byte(a), []byte(b))
}

func formatURLPath(in string) string {
	in = strings.TrimSpace(in)
	if strings.HasSuffix(in, "/") {
		return in[:len(in)-1]
	}
	return in
}

func AKSKAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			ak, sk, sign, timeStamp, serverSign string
			iTime, timeDiff                     int64
			err                                 error
			requestBody                         []byte
		)

		ak = c.Request.Header.Get("AccessKey")
		sign = c.Request.Header.Get("Signature")
		timeStamp = c.Request.Header.Get("TimeStamp")
		if ak == "" || sign == "" || timeStamp == "" {
			abort(c, "header missed: AccessKey|Signature|TimeStamp")
			return
		}

		// check time
		iTime, err = strconv.ParseInt(timeStamp, 10, 64)
		if err != nil {
			abort(c, fmt.Sprintf("TimeStamp Error: %s", err.Error()))
			return
		}
		timeDiff = time.Now().Unix() - iTime
		if timeDiff >= 60 || timeDiff <= -60 {
			abort(c, fmt.Sprintf("timestamp out of range (±60s), diff=%ds", timeDiff))
			return
		}

		// check signature
		sk = getSecKey(ak)
		if sk == "" {
			abort(c, "invalid AccessKey")
			return
		}
		requestBody, err = io.ReadAll(c.Request.Body)
		if err != nil {
			abort(c, err.Error())
			return
		}
		_ = c.Request.Body.Close()
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		serverSign = generateSign(c.Request.Method, formatURLPath(c.Request.URL.Path), c.Request.URL.RawQuery, ak, timeStamp, sk, requestBody)
		if !signEqual(serverSign, sign) {
			ylog.Errorf("AKSKAuth", "signature mismatch for ak=%s method=%s path=%s", ak, c.Request.Method, c.Request.URL.Path)
			abort(c, "signature error")
			return
		}
		c.Next()
	}
}

func abort(c *gin.Context, reason string) {
	c.Abort()
	common.CreateResponse(c, common.AuthFailedErrorCode, reason)
	return
}

func beforeRequestFuncWithKey(req *http.Request, ak, sk string) error {
	var (
		timestamp   = fmt.Sprintf(`%d`, time.Now().Unix())
		err         error
		requestBody []byte
	)

	if req.Body != nil {
		requestBody, err = io.ReadAll(req.Body)
		if err != nil {
			ylog.Errorf("beforeRequestFuncWithKey", "ioutil.ReadAll error %s", err.Error())
			return err
		}
		//Reset after reading
		_ = req.Body.Close()
		req.Body = io.NopCloser(bytes.NewBuffer(requestBody))
	} else {
		requestBody = []byte{}
	}
	sign := generateSign(req.Method, formatURLPath(req.URL.Path), req.URL.RawQuery, ak, timestamp, sk, requestBody)
	req.Header.Add("AccessKey", ak)
	req.Header.Add("Signature", sign)
	req.Header.Add("TimeStamp", timestamp)
	return nil
}

func AuthRequestOption() *grequests.RequestOptions {
	option := &grequests.RequestOptions{
		InsecureSkipVerify: true,
		BeforeRequest:      beforeRequestFunc,
	}
	return option
}

func beforeRequestFunc(req *http.Request) error {
	return beforeRequestFuncWithKey(req, ak, sk)
}
