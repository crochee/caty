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

type People struct {
	Name string
	Age  int
}

func testContent(ctx *gin.Context) {
	Success(ctx, &People{
		Name: "lihua",
		Age:  26,
	}, "分布式网络")
}
