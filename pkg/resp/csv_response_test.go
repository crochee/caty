package resp

import (
	"context"
	"encoding/base64"
	"github.com/minio/minio-go/v7"
	"io"
	"net/http"
	"net/url"
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

func TestName(t *testing.T) {
	data := []byte{
		131, 164, 78, 97, 109, 101, 163, 108, 99, 102, 167, 67, 114, 101, 97, 116, 101, 100, 199, 12, 5, 0, 0, 0, 0, 97, 204, 82, 41, 30, 113, 93, 58, 176, 80, 111, 108, 105, 99, 121, 67, 111, 110, 102, 105, 103, 74, 83, 79, 78, 196, 0,
	}
	t.Logf("%s", data)
	ur := "op/lio/p.log"
	temp := url.QueryEscape(ur)
	t.Log(temp)
	t.Log(base64.URLEncoding.EncodeToString([]byte(temp)))
	c, err := minio.New("", &minio.Options{})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	t.Log(c.MakeBucket(ctx, "", minio.MakeBucketOptions{
		Region:        "",
		ObjectLocking: false,
	}))
}
