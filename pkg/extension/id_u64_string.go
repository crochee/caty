package extension

import (
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/json-iterator/go"
)

func init() {
	jsoniter.ConfigCompatibleWithStandardLibrary.RegisterExtension(&u64AsStringCodec{})
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
							_, _ = stream.Write([]byte(nil))
						} else {
							_, _ = stream.Write([]byte(strconv.FormatUint(val, 10)))
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
