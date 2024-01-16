package models

// 定义上传日志表结构
type UploadLogModel struct {
	ID               uint   `gorm:"primary_key" json:"id"`
	IP               string `gorm:"column:ip" json:"ip"`
	Appkey           string `gorm:"column:appkey" json:"appkey"`
	UploadTime       string `gorm:"column:upload_time" json:"upload_time"`
	FileName         string `gorm:"column:filename" json:"filename"`
	FileSize         string `gorm:"column:file_size" json:"file_size"`
	FileMD5          string `gorm:"column:file_md5" json:"file_md5"`
	FileCreatedTime  int64  `gorm:"column:file_created_time" json:"file_created_time"`
	FileModifiedTime int64  `gorm:"column:file_modified_time" json:"file_modified_time"`
}

// 指定表名
func (UploadLogModel) TableName() string {
	return "t_upload_log"
}
