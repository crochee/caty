package csv

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/afero/mem"
)

func NewCsvRender(file *mem.File, r *http.Request) *csvResponse {
	return &csvResponse{
		file:    file,
		request: r,
	}
}

type csvResponse struct {
	file    *mem.File
	request *http.Request
}

func (c *csvResponse) Render(writer http.ResponseWriter) error {
	now := time.Now().Local()
	fileName := url.QueryEscape(fmt.Sprintf("dcs_%s_%s.csv", c.file.Name(), now.Format("2006-01-02")))
	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	http.ServeContent(writer, c.request, fileName, c.file.Info().ModTime(), c.file)
	return nil
}

func (c *csvResponse) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
}
