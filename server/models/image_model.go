package models

import "gorm.io/gorm"

type UploadImageModel struct {
	gorm.Model             // gorm.Model 包含了一些常见的字段，例如 ID、CreatedAt、UpdatedAt、DeletedAt 等，这些字段通常用于数据库记录的标识和时间戳。
	FileUUID        string `gorm:"column:file_uuid"`
	Size            int64  `gorm:"column:size"`
	FileName        string `gorm:"column:file_name"`
	FilePath        string `gorm:"column:file_path"`
	ThumbFilePath   string `gorm:"column:thumb_file_path"`
	ThumbnailBase64 string `gorm:"column:thumbnail_base64"`
	FileMD5         string `gorm:"column:file_md5"`
	DownloadCount   int    `gorm:"column:download_count"`
	Sort            int    `gorm:"column:sort"`
	UploadTime      string `gorm:"column:upload_time" json:"-"`
	UploadIP        string `gorm:"column:upload_ip" json:"-"`
	UploadUser      string `gorm:"column:upload_user" json:"-"`
	Status          string `gorm:"column:status"`
}

func (UploadImageModel) TableName() string {
	return "t_upload_image"
}
