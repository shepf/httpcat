package models

import "time"

// ShareModel 文件分享表
type ShareModel struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	ShareCode    string    `gorm:"column:share_code;uniqueIndex;size:16" json:"shareCode"`   // 分享短码
	FilePath     string    `gorm:"column:file_path;not null" json:"filePath"`                // 文件相对路径
	FileName     string    `gorm:"column:file_name;not null" json:"fileName"`                // 文件名（展示用）
	FileType     string    `gorm:"column:file_type;default:file" json:"fileType"`            // file / image
	CreatedBy    string    `gorm:"column:created_by;not null" json:"createdBy"`              // 创建者
	ExtractCode  string    `gorm:"column:extract_code" json:"extractCode"`                   // 提取码（空=无需提取码）
	ExpireAt     *time.Time `gorm:"column:expire_at" json:"expireAt"`                        // 过期时间（NULL=永不过期）
	MaxDownloads int       `gorm:"column:max_downloads;default:0" json:"maxDownloads"`       // 最大下载次数（0=不限）
	CurDownloads int       `gorm:"column:cur_downloads;default:0" json:"curDownloads"`       // 当前已下载次数
	IsActive     bool      `gorm:"column:is_active;default:true" json:"isActive"`            // 是否有效
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (ShareModel) TableName() string {
	return "t_share"
}
