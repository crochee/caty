// Date: 2021/9/19

// Package validator
package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/zh"
)

// Init init validator
func Init() error {
	v := &defaultValidator{validate: validator.New()}
	v.validate.SetTagName("binding")
	v.translator, _ = ut.New(zh.New()).GetTranslator("zh")
	if err := translations.RegisterDefaultTranslations(v.validate, v.translator); err != nil {
		return err
	}
	binding.Validator = v
	return nil
}

type defaultValidator struct {
	validate   *validator.Validate
	translator ut.Translator
}

type sliceValidateError []error

func (err sliceValidateError) Error() string {
	var errs []string
	for i, e := range err {
		if e == nil {
			continue
		}
		errs = append(errs, fmt.Sprintf("[%d]: %s", i, e.Error()))
	}
	return strings.Join(errs, ";")
}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	err := v.defaultValidateStruct(obj)
	if err == nil {
		return nil
	}
	vErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}
	var errs sliceValidateError
	for _, s := range vErrs.Translate(v.translator) {
		errs = append(errs, errors.New(s))
	}
	return errs
}

// validateStruct receives struct type
func (v *defaultValidator) validateStruct(obj interface{}) error {
	return v.validate.Struct(obj)
}

// Engine returns the underlying validator engine which powers the default
// Validator instance. This is useful if you want to register custom validations
// or struct level validations. See validator GoDoc for more info -
// https://godoc.org/gopkg.in/go-playground/validator.v8
func (v *defaultValidator) Engine() interface{} {
	return v.validate
}

func (v *defaultValidator) defaultValidateStruct(obj interface{}) error {
	if obj == nil {
		return nil
	}
	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		return v.ValidateStruct(value.Elem().Interface())
	case reflect.Struct:
		return v.validateStruct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(sliceValidateError, 0)
		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}
