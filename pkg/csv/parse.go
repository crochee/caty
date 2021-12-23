package csv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/crochee/lirity/e"
	"github.com/crochee/lirity/variable"
)

type parse struct {
	tagName string
}

func (p *parse) parseStruct(obj interface{}) ([]map[string]interface{}, error) {
	if obj == nil {
		return nil, errors.New("it's nil")
	}
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, errors.New("not struct")
	}
	temp := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		fv := v.Field(i)
		ft := t.Field(i)

		if !fv.IsValid() {
			continue
		}
		if !fv.CanInterface() {
			continue
		}
		if ft.PkgPath != "" { // unexported
			continue
		}
		name, option := parseTag(ft.Tag.Get(p.tagName))
		if name == "-" {
			continue // ignore "-"
		}
		if name == "" {
			name = ft.Name // use field name
		}
		if option == "omitempty" && fv.IsZero() {
			continue // skip empty field
		}
		// ft.Anonymous means embedded field
		if ft.Anonymous {
			if !fv.IsValid() || fv.IsNil() {
				continue
			}
			embedded, err := p.Parse(fv.Interface())
			if err != nil {
				return nil, err
			}
			for _, embMap := range embedded {
				for embName, embValue := range embMap {
					temp[embName] = embValue
				}
			}
			continue
		}

		if option == "string" {
			tempString := num2String(fv)
			if tempString != nil {
				temp[name] = tempString
				continue
			}
		}
		temp[name] = fv.Interface()
	}
	return []map[string]interface{}{temp}, nil
}

func (p *parse) Parse(obj interface{}) ([]map[string]interface{}, error) {
	if obj == nil {
		return nil, errors.New("it's' nil")
	}
	value := reflect.ValueOf(obj)
	switch value.Kind() { // nolint:exhaustive
	case reflect.Ptr:
		return p.parseStruct(value.Elem().Interface())
	case reflect.Struct:
		return p.parseStruct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(e.Errors, 0, count)
		tempMap := make([]map[string]interface{}, 0, count)
		for i := 0; i < count; i++ {
			if v, err := p.parseStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			} else {
				tempMap = append(tempMap, v...)
			}
		}
		if len(validateRet) == 0 {
			return tempMap, nil
		}
		return nil, validateRet
	default:
		return nil, fmt.Errorf("not support %s", value.Kind().String())
	}
}

func num2String(fv reflect.Value) interface{} {
	kind := fv.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(fv.Int(), variable.DecimalSystem)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(fv.Uint(), variable.DecimalSystem)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(fv.Float(), 'f', 2, 64)
	default:
		panic(fmt.Sprintf("not support %s to string", kind.String()))
	}
}

func parseTag(tag string) (tag0, tag1 string) {
	tags := strings.Split(tag, ",")
	if len(tags) == 0 {
		return
	}
	if len(tags) == 1 {
		tag0 = tags[0]
		return
	}
	tag0, tag1 = tags[0], tags[1]
	return
}
