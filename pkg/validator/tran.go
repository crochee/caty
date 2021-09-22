// Date: 2021/9/19

// Package validator
package validator

import (
	"errors"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/zh"

	"obs/pkg/e"
)

var v *defaultValidator

func Var(field interface{}, tag string) error {
	if v == nil {
		return nil
	}
	err := v.Validate.Var(field, tag)
	if err == nil {
		return nil
	}
	return v.Translate(err)
}

// Init init validator
func Init() error {
	v = &defaultValidator{Validate: validator.New()}
	v.Validate.SetTagName("binding")
	v.translator, _ = ut.New(zh.New()).GetTranslator("zh")
	if err := translations.RegisterDefaultTranslations(v.Validate, v.translator); err != nil {
		return err
	}
	binding.Validator = v
	return nil
}

type defaultValidator struct {
	Validate   *validator.Validate
	translator ut.Translator
}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	err := v.defaultValidateStruct(obj)
	if err == nil {
		return nil
	}
	return v.Translate(err)
}

// Translate receives struct type
func (v *defaultValidator) Translate(err error) error {
	if err == nil {
		return nil
	}
	vErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}
	var errs e.Errors
	for _, s := range vErrs.Translate(v.translator) {
		errs = append(errs, errors.New(s))
	}
	return errs
}

// validateStruct receives struct type
func (v *defaultValidator) validateStruct(obj interface{}) error {
	return v.Validate.Struct(obj)
}

// Engine returns the underlying validator engine which powers the default
// Validator instance. This is useful if you want to register custom validations
// or struct level validations. See validator GoDoc for more info -
// https://godoc.org/gopkg.in/go-playground/validator.v8
func (v *defaultValidator) Engine() interface{} {
	return v.Validate
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
		validateRet := make(e.Errors, 0, count)
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
