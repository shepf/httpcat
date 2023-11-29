package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"gin_web_demo/server/common/ylog"
	"strings"
)

const (
	AuthorizationPrefixHttpCat = "HttpCat "
)

// 凭证 AK/SK可以从  todo 暂时未提供web界面，暂时配置文件配置方式获取
type Credentials struct {
	AccessKey string
	SecretKey []byte
}

// 构建一个Credentials对象
func New(accessKey, secretKey string) *Credentials {
	return &Credentials{accessKey, []byte(secretKey)}
}

// Sign 对数据进行签名，一般用于私有空间下载用途
// Sign 方法使用 HMAC 和 SHA-1 算法对给定的数据进行签名，使用访问密钥中的 Secret Key 作为密钥。
// 最终的签名结果包含了 Access Key 和 HMAC 的摘要，构成了一个格式化的字符串。
func (ath *Credentials) Sign(data []byte) (token string) {
	h := hmac.New(sha1.New, ath.SecretKey)
	h.Write(data)

	sign := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s:%s", ath.AccessKey, sign)
}

// SignWithData 对数据进行签名，一般用于上传凭证的生成用途
func (ath *Credentials) SignWithData(b []byte) (token string) {
	encodedData := base64.URLEncoding.EncodeToString(b)
	sign := ath.Sign([]byte(encodedData))
	return fmt.Sprintf("%s:%s", sign, encodedData)
}

func (ath *Credentials) VerifyUploadToken(token string) bool {
	parts := strings.Split(token, ":")
	if len(parts) != 3 {
		ylog.Errorf("[AUTH]", "token 格式错误")
		return false
	}

	accessKey := parts[0]
	signature := parts[1]
	putPolicyBase64 := parts[2]

	if accessKey != ath.AccessKey {
		fmt.Println(accessKey)
		fmt.Println(ath.AccessKey)
		ylog.Errorf("[AUTH]", "accessKey: %s", ath.AccessKey)
		ylog.Errorf("[AUTH]", "ath.AccessKey: %s", ath.AccessKey)

		ylog.Errorf("[AUTH]", "accessKey 不匹配")
		return false
	}

	// 解析上传策略
	putPolicyJSON, err := base64.StdEncoding.DecodeString(putPolicyBase64)
	if err != nil {
		ylog.Errorf("[AUTH]", "上传策略解析失败")
		return false
	}

	// 使用 AccessKey 对上传策略进行签名
	expectedSignature := ath.SignWithData(putPolicyJSON)
	ylog.Infof("[AUTH]", "expectedSignature: %s", expectedSignature)
	ylog.Infof("[AUTH]", "signature: %s", signature)
	ylog.Infof("[AUTH]", "token: %s", token)
	// 比较签名
	return token == expectedSignature
}
