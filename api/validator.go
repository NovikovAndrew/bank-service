package api

import (
	"bank-service/util"
	"github.com/go-playground/validator/v10"
)

var currencyValidator validator.Func = func(fl validator.FieldLevel) bool {
	currency, ok := fl.Field().Interface().(string)

	if ok {
		isFind := false
		for _, c := range util.Currencies {
			if c == currency {
				isFind = true
				break
			}
		}
		return isFind
	}

	return false
}
