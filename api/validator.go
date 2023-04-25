package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/vutranhoang1411/SimpleBank/util"
)
var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency,ok:=fl.Field().Interface().(string);ok{
		return util.CheckCurrency(currency)
	}
	return false;
}
var validEmail validator.Func=func(fl validator.FieldLevel) bool {
	if email,ok:=fl.Field().Interface().(string);ok{
		return util.CheckEmail(email);
	}
	return false;
}