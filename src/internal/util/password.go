package util

import (
	"regexp"

	"github.com/go-playground/validator"
)

func PasswordValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		hasUpper  = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower  = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSymbol = regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)
	)

	return hasUpper && hasLower && hasNumber && hasSymbol
}
