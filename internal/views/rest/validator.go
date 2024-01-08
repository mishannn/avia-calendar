package rest

import (
	"github.com/go-playground/validator"
)

type RequestValidator struct {
	validator *validator.Validate
}

func (rv *RequestValidator) Validate(i interface{}) error {
	return rv.validator.Struct(i)
}
