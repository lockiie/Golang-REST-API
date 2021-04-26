package models

import (
	"eco/src/types"
	"errors"

	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	pt_translations "github.com/go-playground/validator/v10/translations/pt_BR"
)

var (
	validate *validator.Validate
	uni      *ut.UniversalTranslator
	trans    ut.Translator
)

func init() {
	validate = validator.New()
	pt := pt_BR.New()
	uni = ut.New(pt, pt)
	trans, _ = uni.GetTranslator("pt_BR")
	validate = validator.New()
	pt_translations.RegisterDefaultTranslations(validate, trans)
}

func translateError(err error) error {
	if err == nil {
		return nil
	}
	validatorErrs := err.(validator.ValidationErrors)
	var message string

	for _, e := range validatorErrs {
		if message != types.EmptyStr {
			message += ", "
		}
		message += e.Translate(trans)
	}
	return errors.New(message)
}
