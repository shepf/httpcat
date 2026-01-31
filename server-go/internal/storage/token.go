package storage

import (
	"encoding/json"
	"httpcat/internal/storage/auth"
)

type UploadPolicy struct {

	// 上传凭证有效截止时间。Unix时间戳，单位为秒。该截止时间为上传完成后，在httpcat生成文件的校验时间，而非上传的开始时间，
	// 一般建议设置为上传开始时间 + 3600s，用户可根据具体的业务场景对凭证截止时间进行调整。
	Deadline uint64 `json:"deadline"`

	// 限定上传文件大小最小值，单位Byte。小于限制上传文件大小的最小值会被判为上传失败，返回 403 状态码
	FsizeMin int64 `json:"fsizeMin,omitempty"`

	// 限定上传文件大小最大值，单位Byte。超过限制上传文件大小的最大值会被判为上传失败，返回 413 状态码。
	FsizeLimit int64 `json:"fsizeLimit,omitempty"`

	// 接收持久化处理结果通知的 URL。必须是公网上可以正常进行 POST 请求并能响应 HTTP/1.1 200 OK 的有效 URL。该 URL 获取的内容和持久化处
	// 理状态查询的处理结果一致。发送 body 格式是 Content-Type 为 application/json 的 POST 请求，需要按照读取流的形式读取请求的 body
	// 才能获取。
	PersistentNotifyURL string `json:"persistentNotifyUrl,omitempty"`
}

// UploadToken 方法用来进行上传凭证的生成
// 该方法生成的过期时间是现对于现在的时间
func (p *UploadPolicy) UploadToken(cred *auth.Credentials) string {
	return p.uploadToken(cred)
}

func (p UploadPolicy) uploadToken(cred *auth.Credentials) (token string) {
	// 暂时注释上传策略使用全局暂时不和用户关联，这样生成的上传token一样
	//if p.Deadline == 0 {
	//	p.Deadline = 7200 // 默认一小时过期
	//}
	//p.Deadline += uint64(time.Now().Unix())
	putPolicyJSON, _ := json.Marshal(p)
	token = cred.SignWithData(putPolicyJSON)
	return
}
