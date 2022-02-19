package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"

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

	if _, err = m.w.Write([]byte("\xEF\xBB\xBF")); err != nil { // 写入UTF-8 BOM，防止中文乱码
		return err
	}
	w := csv.NewWriter(m.w)

	header := make([]string, 0, len(headers))
	for _, key := range headers {
		header = append(header, key.key)
	}
	if err = w.Write(header); err != nil {
		return err
	}

	rows := make([][]string, len(data))
	for i, d := range data {
		row := make([]string, len(header))
		for index, key := range header {
			v, found := d.data[key]
			if found {
				row[index] = toString(v)
			}
		}
		rows[i] = row
	}
	return w.WriteAll(rows)
}

func toString(d interface{}) string {
	switch v := d.(type) {
	case bool:
		if v {
			return "true"
		}
		return "false"
	case *bool:
		if *v {
			return "true"
		}
		return "false"
	case string:
		return v
	case *string:
		return *v
	case int, *int, int64, *int64:
		return fmt.Sprintf("%d", d)
	default:
		j, _ := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(d)
		return j
	}
}
