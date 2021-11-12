// Package model
package model

import (
	"gorm.io/gorm"

	"caty/pkg/v"
)

type Page struct {
	// 分页索引
	Index uint64 `json:"num" form:"index"`
	// 分页大小
	Size int `json:"size" form:"size"`
	// 总数
	Total int `json:"total" form:"total"`
}

func HandlePage(query *gorm.DB, page Page) *gorm.DB {
	if page.Size == v.PageAll {
		return query
	} else if page.Size == 0 {
		if page.Index == 0 {
			return query
		}
		size := v.DefaultPageSize
		index := int(page.Index)
		return query.Limit(size).Offset((index - 1) * size)
	} else if page.Size > 0 {
		var index int
		if page.Index == 0 {
			index = v.DefaultPageIndex
		} else {
			index = int(page.Index)
		}
		return query.Limit(page.Size).Offset((index - 1) * page.Size)
	}
	return query
}
