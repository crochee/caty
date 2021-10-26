// Package model
package model

type Page struct {
	// 分页索引
	Index int `json:"num" form:"index"`
	// 分页大小
	Size int `json:"size" form:"size"`
	// 总数
	Total int `json:"total" form:"total"`
}

const (
	DefaultIndex = 1
	DefaultSize  = 20
)

func DefaultPage() *Page {
	return &Page{
		Index: DefaultIndex,
		Size:  DefaultSize,
	}
}
