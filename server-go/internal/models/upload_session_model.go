package models

import "time"

// UploadSessionModel 分片上传会话（v0.7.0 新增）
// 用于支持大文件分片上传 + 断点续传：
//   1. 客户端调用 /api/v1/file/upload/init 创建会话，拿到 UploadID
//   2. 客户端并发上传分片 /api/v1/file/upload/chunk，服务端记录已上传分片索引
//   3. 所有分片上传完成后，调用 /api/v1/file/upload/complete 合并
//   4. 任意时刻调用 /api/v1/file/upload/status 查询已上传分片，用于断点续传
type UploadSessionModel struct {
	ID           uint   `gorm:"primarykey" json:"id"`
	UploadID     string `gorm:"column:upload_id;uniqueIndex;size:64;not null" json:"uploadId"` // 会话唯一 ID（UUID）
	Appkey       string `gorm:"column:appkey;size:64" json:"appkey"`                           // 上传凭证 appkey
	Username     string `gorm:"column:username;size:64;index" json:"username"`                 // 发起用户
	FileName     string `gorm:"column:file_name;not null" json:"fileName"`                     // 最终文件名
	RelDir       string `gorm:"column:rel_dir;size:1024" json:"relDir"`                        // 相对于上传根的子目录（空=根目录）
	FileSize     int64  `gorm:"column:file_size;not null" json:"fileSize"`                     // 总文件大小（字节）
	ChunkSize    int64  `gorm:"column:chunk_size;not null" json:"chunkSize"`                   // 每片大小（最后一片可小于）
	TotalChunks  int    `gorm:"column:total_chunks;not null" json:"totalChunks"`               // 总分片数
	FileMD5      string `gorm:"column:file_md5;size:64;index" json:"fileMD5"`                  // 客户端声明的整体 MD5（可选，用于秒传/校验）
	UploadedBits string `gorm:"column:uploaded_bits;type:text" json:"uploadedBits"`            // 已上传分片 bitmap（"01011..." 字符串，每位对应一个分片）
	UploadedNum  int    `gorm:"column:uploaded_num;default:0" json:"uploadedNum"`              // 已上传分片数（缓存，避免扫描 bitmap）
	Status       string `gorm:"column:status;size:16;index;default:active" json:"status"`      // active / completed / aborted
	FinalPath    string `gorm:"column:final_path;size:1024" json:"finalPath"`                  // 完成后的最终文件绝对路径
	IP           string `gorm:"column:ip;size:64" json:"ip"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updatedAt"`
	ExpireAt     time.Time `gorm:"column:expire_at;index" json:"expireAt"` // 会话过期时间（默认 24h，过期后分片目录会被清理）
}

func (UploadSessionModel) TableName() string {
	return "t_upload_session"
}
