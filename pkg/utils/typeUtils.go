package utils

import (
	"cloudProject/pkg/cast"
	"reflect"
)

func IsNonEmptyString(i interface{}) bool {
	variable := reflect.ValueOf(i)
	msg := ""
	if variable.Kind().String() == "string" {
		msg, _ = cast.ToString(i)
	}
	return "" == msg
}
func IsBool(i interface{}) bool {
	variable := reflect.ValueOf(i)
	return variable.Kind().String() == "bool"
}
