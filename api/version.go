// Date: 2021/10/15

// Package api
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type VersionResponse struct {
	// 结果集
	// Required: true
	Result []*VersionResult `json:"result"`
}

const (
	PreRelease = "pre-release"
	Online     = "online"
	Offline    = "offline"
)

type VersionResult struct {
	// 版本
	// Required: true
	Version string `json:"version"`
	// 状态 pre-release,online,offline
	// Required: true
	Status string `json:"status"`
	// 发布时间
	// Required: true
	Release string `json:"release"`
	// 下线时间
	Offline string `json:"offline"`
}

// Version godoc
// swagger:operation GET / 通用 SNullRequest
// ---
// summary: 查询api版本信息
// description: 查询api版本详细信息
// produces:
// - application/json
// responses:
//   '200':
//     type: object
//     "$ref": "#/responses/SAPIVersionResponse"
func Version(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, &VersionResponse{Result: []*VersionResult{
		{
			Version: "v1",
			Status:  Online,
			Release: "2021-10- 15:04:05",
			Offline: "",
		},
	}})
}
