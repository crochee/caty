package resp

import (
	"io"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"

	"caty/internal"
)

func TestMarshal(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/files", testContent)
	headers := make(http.Header)
	headers.Set("Accept", "text/csv")
	w := internal.PerformRequest(router, http.MethodGet, "/files", nil, headers).Result()
	defer w.Body.Close()
	data, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d %v\n%s", w.StatusCode, w.Header, data)
}

type Content struct {
	Index int    `csv:"index,5"`
	Some  string `csv:",6"`
}
type People struct {
	Name     string     `csv:",1"`
	Age      int        `csv:",2"`
	Index    int        `csv:",3"`
	Create   int        `csv:""`
	Contents []*Content `csv:"文章,fmt"`
}

func testContent(ctx *gin.Context) {
	Success(ctx, &People{
		Name:   "lihua",
		Age:    26,
		Index:  9,
		Create: 0,
		Contents: []*Content{
			{
				Index: 0,
				Some:  "t",
			},
			{
				Index: 1,
				Some:  "s",
			},
			{
				Index: 2,
				Some:  "p",
			},
		},
	}, "分布式网络")
}
