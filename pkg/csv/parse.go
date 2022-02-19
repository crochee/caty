package csv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/crochee/lirity/e"
	"github.com/crochee/lirity/variable"
	"github.com/json-iterator/go"
)

type mapIndexValue struct {
	data  map[string]interface{}
	index []*indexValue
}

type indexValue struct {
	key   string
	index uint64
}

type parse struct {
	tagName string
}

func (p *parse) parseStruct(obj interface{}) ([]*mapIndexValue, error) {
	if obj == nil {
		return []*mapIndexValue{}, nil
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
	result := &mapIndexValue{
		data:  map[string]interface{}{},
		index: []*indexValue{},
	}

	for i := 0; i < t.NumField(); i++ {
		fv := v.Field(i)
		ft := t.Field(i)

		if !fv.IsValid() {
			continue
		}
		if !fv.CanInterface() {
			continue
		}
		if !ft.IsExported() { // unexported
			continue
		}
		tag, found := ft.Tag.Lookup(p.tagName)
		if !found {
			continue
		}
		tags := strings.Split(tag, ",")
		var (
			name   string
			option string
			index  int64
		)
		if len(tags) == 1 {
			name = tags[0]
		} else if len(tags) == 2 {
			name = tags[0]
			if tags[1] == "string" || tags[1] == "fmt" {
				option = tags[1]
			} else {
				index, _ = strconv.ParseInt(tags[1], variable.DecimalSystem, 64)
			}
		} else {
			name = tags[0]
			option = tags[1]
			index, _ = strconv.ParseInt(tags[2], variable.DecimalSystem, 64)
		}
		if name == "-" {
			continue // ignore "-"
		}
		if name == "" {
			name = ft.Name // use field name
		}
		if ft.Anonymous || fv.Kind() == reflect.Slice || fv.Kind() == reflect.Array ||
			fv.Kind() == reflect.Struct || fv.Kind() == reflect.Ptr {
			if fv.IsZero() {
				continue
			}
			embedded, err := p.parse(fv.Interface())
			if err != nil {
				return nil, err
			}
			if (fv.Kind() == reflect.Slice || fv.Kind() == reflect.Array) && option == "fmt" {
				// fmt
				embedded = format(name, embedded)
			}

			for _, embMap := range embedded {
				for embName, embValue := range embMap.data {
					result.data[embName] = embValue
				}
				for _, embIndexValue := range embMap.index {
					result.index = append(result.index, &indexValue{
						key:   embIndexValue.key,
						index: embIndexValue.index,
					})
				}
			}
			continue
		}
		if option == "string" {
			tempString := value2String(fv)
			if tempString != nil {
				result.data[name] = tempString
				result.index = append(result.index, &indexValue{
					key:   name,
					index: uint64(index),
				})
				continue
			}
		}
		result.data[name] = fv.Interface()
		result.index = append(result.index, &indexValue{
			key:   name,
			index: uint64(index),
		})
	}
	return []*mapIndexValue{result}, nil
}

func (p *parse) parse(obj interface{}) ([]*mapIndexValue, error) {
	if obj == nil {
		return []*mapIndexValue{}, nil
	}
	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		value = value.Elem()
		if !value.IsValid() {
			return []*mapIndexValue{}, nil
		}
		return p.parse(value.Interface())
	case reflect.Struct:
		return p.parseStruct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(e.Errors, 0, count)
		tempMap := make([]*mapIndexValue, 0, count)
		for i := 0; i < count; i++ {
			if !value.Index(i).CanInterface() {
				validateRet = append(validateRet, fmt.Errorf("%s can't interface", value.Index(i).String()))
				continue
			}
			if v, err := p.parse(value.Index(i).Interface()); err != nil {
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

func value2String(fv reflect.Value) interface{} {
	kind := fv.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(fv.Int(), variable.DecimalSystem)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(fv.Uint(), variable.DecimalSystem)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(fv.Float(), 'f', 2, 64)
	default:
		data, _ := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(fv.Interface())
		return data
	}
}

type indexValueList []*indexValue

func (l indexValueList) Len() int {
	return len(l)
}

func (l indexValueList) Less(i, j int) bool {
	return l[i].index < l[j].index
}

func (l indexValueList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func format(name string, input []*mapIndexValue) []*mapIndexValue {
	result := &mapIndexValue{
		data:  map[string]interface{}{},
		index: make([]*indexValue, 0, 4),
	}
	for _, embMap := range input {
		for embName, embValue := range embMap.data {
			embName = fmt.Sprintf("%s(%s)", name, embName)
			v, ok := result.data[embName]
			if !ok {
				result.data[embName] = embValue
				continue
			}
			result.data[embName] = fmt.Sprintf("%v,%v", v, embValue)
		}
		for _, embIndexValue := range embMap.index {
			var ignore bool
			embName := fmt.Sprintf("%s(%s)", name, embIndexValue.key)
			for _, v := range result.index {
				if v.key == embName {
					ignore = true
					break
				}
			}
			if !ignore {
				result.index = append(result.index, &indexValue{
					key:   embName,
					index: embIndexValue.index,
				})
			}
		}
	}
	return []*mapIndexValue{result}
}
