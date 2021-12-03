package validate

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Struct(s interface{}) (string, error) {
	err := validate.Struct(s)
	if err != nil {
		customError := err.(validator.ValidationErrors)
		return customError[0].StructField() + "," + customError[0].ActualTag(), err
	}
	return "", nil
}
