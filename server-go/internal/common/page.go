package common

const (
	DefaultPage     = 1
	DefaultPageSize = 100
)

type PageRequest struct {
	Page       int64  `form:"page,default=1" binding:"required,numeric,min=1"`
	PageSize   int64  `form:"page_size,default=100" binding:"required,numeric,min=1,max=5000"`
	OrderKey   string `form:"order_key"`
	OrderValue int    `form:"order_value"`
}

type PageOption struct {
	Page     int64
	PageSize int64
	Filter   interface{}
	Sorter   interface{}
}

// PageSearch 查询定义
type PageSearch struct {
	Page     int64
	PageSize int64
	Filter   interface{}
	Sorter   interface{}
}

// PageResponse 返回定义
type PageResponse struct {
	Total    int64 `json:"total" bson:"total"`
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
}

type ModelPage struct {
	Total    int64         `json:"total"`
	Count    int64         `json:"count"`
	Pages    int64         `json:"pages"`
	Page     int64         `json:"page"`
	PageSize int64         `json:"page_size"`
	HasPrev  bool          `json:"has_prev"`
	HasNext  bool          `json:"has_next"`
	Items    []interface{} `json:"items"`
}
