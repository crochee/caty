package validator

import (
	"github.com/crochee/lirity/validator"
	"github.com/gin-gonic/gin/binding"
)

// Init init validator
func Init() error {
	v, err := validator.New()
	if err != nil {
		return err
	}
	binding.Validator = v
	return nil
}
