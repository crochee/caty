package internal

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"unsafe"

	"github.com/json-iterator/go"
)

type TestString struct {
	ID     uint64 `json:"id,string"`
	PID    uint64 `json:"pid"`
	PartID uint64 `json:",string"`
}

func TestStringName(t *testing.T) {
	jsoniter.ConfigCompatibleWithStandardLibrary.RegisterExtension(&u64AsStringCodec{})
	data, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(&TestString{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", data)
	if data, err = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(&TestString{
		ID: 787446465166,
	}); err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", data)
	if data, err = json.Marshal(&TestString{
		ID: 787446465166,
	}); err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", data)
	var t1 TestString
	originData := []byte(`{"id":"","pid":0,"PartID":""}`)
	if err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(originData, &t1); err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", t1)
	var t2 TestString
	originData = []byte(`{"id":"787446465166","pid":0,"PartID":""}`)
	if err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(originData, &t2); err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", t2)
	var t3 TestString
	originData = []byte(`{"id":"787446465166","pid":0,"PartID":"0"}`)
	if err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(originData, &t3); err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", t3)
}

type u64AsStringCodec struct {
	jsoniter.DummyExtension
}

func (extension *u64AsStringCodec) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for _, binding := range structDescriptor.Fields {
		if binding.Field.Type().Kind() == reflect.Uint64 {
			tagParts := strings.Split(binding.Field.Tag().Get("json"), ",")
			if len(tagParts) <= 1 {
				continue
			}
			for _, tagPart := range tagParts[1:] {
				if tagPart == "string" {
					binding.Encoder = &funcEncoder{fun: func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
						val := *((*uint64)(ptr))
						if val == 0 {
							stream.Write([]byte(nil))
						} else {
							stream.Write([]byte(strconv.FormatUint(val, 10)))
						}
					}}
					binding.Decoder = &funcDecoder{func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
						if iter.WhatIsNext() != jsoniter.StringValue {
							*((*uint64)(ptr)) = iter.ReadUint64()
						}
					}}
					break
				}
			}
		}
	}
}

type funcEncoder struct {
	fun         jsoniter.EncoderFunc
	isEmptyFunc func(ptr unsafe.Pointer) bool
}

func (encoder *funcEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	encoder.fun(ptr, stream)
}

func (encoder *funcEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	if encoder.isEmptyFunc == nil {
		return false
	}
	return encoder.isEmptyFunc(ptr)
}

type funcDecoder struct {
	fun jsoniter.DecoderFunc
}

func (decoder *funcDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	decoder.fun(ptr, iter)
}
