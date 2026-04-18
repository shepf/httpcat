package models

// OperationLogModel 操作日志表
type OperationLogModel struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Username  string `gorm:"column:username;index" json:"username"`   // 操作用户
	IP        string `gorm:"column:ip" json:"ip"`                     // 请求 IP
	Method    string `gorm:"column:method" json:"method"`             // HTTP 方法
	Path      string `gorm:"column:path" json:"path"`                 // 请求路径
	Action    string `gorm:"column:action;index" json:"action"`       // 操作类型：upload/download/delete/rename/mkdir/share/preview/login/config 等
	Detail    string `gorm:"column:detail;size:1024" json:"detail"`   // 操作详情（如文件名、目标路径等）
	Status    int    `gorm:"column:status" json:"status"`             // HTTP 状态码
	Latency   int64  `gorm:"column:latency" json:"latency"`          // 请求耗时（毫秒）
	UserAgent string `gorm:"column:user_agent;size:512" json:"userAgent"` // 用户代理
	CreatedAt string `gorm:"column:created_at;index" json:"createdAt"`    // 操作时间
}

func (OperationLogModel) TableName() string {
	return "t_operation_log"
}
