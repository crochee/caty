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
	Writer     io.Writer
}

func NewMarshal(opts ...func(*Option)) *marshal {
	option := Option{
		TagName:    "csv",
		FieldNames: []string{"Result", "List"},
		Writer:     os.Stdout,
	}
	for _, opt := range opts {
		opt(&option)
	}
	return &marshal{
		fieldNames: option.FieldNames,
		expose:     expose{},
		parser:     &parse{tagName: option.TagName},
		w:          option.Writer,
	}
}

type marshal struct {
	fieldNames []string
	expose     expose
	parser     *parse
	w          io.Writer
}

func (m *marshal) Encode(obj interface{}) error {
	value, err := m.expose.GetStruct(obj, m.fieldNames...)
	if err != nil {
		return err
	}
	var data []*mapIndexValue
	if data, err = m.parser.parse(value); err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	headers := make([]*indexValue, 0, 4)
	for _, v := range data {
		for _, indexKey := range v.index {
			var ignore bool
			for _, header := range headers {
				if indexKey.key == header.key {
					ignore = true
					break
				}
			}
			if !ignore {
				headers = append(headers, &indexValue{
					key:   indexKey.key,
					index: indexKey.index,
				})
			}
		}
	}
	sort.Sort(indexValueList(headers))
	header := make(table.Row, 0, len(data[0].index))
	for _, key := range headers {
		header = append(header, key.key)
	}

	rows := make([]table.Row, len(data))
	for i, d := range data {
		row := make(table.Row, len(headers))
		for index, key := range headers {
			v, found := d.data[key.key]
			if found {
				row[index] = text.WrapHard(toString(v), DefaultTransverseStringLength)
			}
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
