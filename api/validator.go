package api

import (
	"simplebank/util"

	"github.com/go-playground/validator/v10"
)

// making validator for type of currency name

var validCurrency validator.Func = func(FieldLevel validator.FieldLevel) bool {
	if currency, ok := FieldLevel.Field().Interface().(string); ok {
		// check if currency is supported or not
		return util.IsSupportedCurrency(currency)
	}
	return false

}
