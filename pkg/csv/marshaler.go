package csv

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/json-iterator/go"
)

type Option struct {
	TagName    string
	FieldNames []string
}

func NewMarshal(opts ...func(*Option)) *marshal {
	option := Option{TagName: "csv", FieldNames: []string{"Result", "List"}}
	for _, opt := range opts {
		opt(&option)
	}
	return &marshal{
		fieldNames: option.FieldNames,
		expose:     expose{},
		parser:     &parse{tagName: option.TagName},
		w:          os.Stdout,
	}
}

type marshal struct {
	fieldNames []string
	expose     expose
	parser     *parse
	w          io.Writer
}

func (m *marshal) NewEncoder(w io.Writer) {
	m.w = w
}

func (m *marshal) Encode(obj interface{}) error {
	value, err := m.expose.GetStruct(obj, m.fieldNames...)
	if err != nil {
		return err
	}
	var data []map[string]interface{}
	if data, err = m.parser.Parse(value); err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	headerStrings := make([]string, 0, len(data[0]))
	for key := range data[0] {
		headerStrings = append(headerStrings, key)
	}
	sort.Strings(headerStrings)
	header := make(table.Row, 0, len(headerStrings))
	for _, headerValue := range headerStrings {
		header = append(header, headerValue)
	}

	rows := make([]table.Row, len(data))
	for i, d := range data {
		row := make(table.Row, len(headerStrings))
		for k, v := range d {
			index := indexOf(headerStrings, k)
			if index == -1 {
				continue
			}
			row[index] = text.WrapHard(toString(v), DefaultTransverseStringLength)
		}
		rows[i] = row
	}

	t := table.NewWriter()
	t.SetOutputMirror(m.w)
	t.AppendHeader(header)
	t.AppendRows(rows)
	t.RenderCSV()
	return nil
}

const (
	DefaultTransverseStringLength = 64
)

func indexOf(list []string, target string) int {
	for i, v := range list {
		if v == target {
			return i
		}
	}
	return -1
}

func toString(d interface{}) string {
	switch d.(type) {
	case bool, *bool:
		return fmt.Sprintf("%t", d)
	case string, *string:
		return fmt.Sprintf("%s", d)
	case int, *int, int64, *int64:
		return fmt.Sprintf("%d", d)
	default:
		j, _ := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(d)
		return j
	}
}
