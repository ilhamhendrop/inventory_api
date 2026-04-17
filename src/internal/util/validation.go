package util

import (
	"fmt"

	"github.com/go-playground/validator"
)

func Validate[T any](data T) string {
	validate := validator.New()
	validate.RegisterValidation("password", PasswordValidation)
	err := validate.Struct(data)

	if err == nil {
		return ""
	}

	errMap := map[string]string{}

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, v := range errs {
			errMap[v.Field()] = TranslateTag(v)
		}
	}

	msg := ""
	for _, v := range errMap {
		if msg != "" {
			msg += ", "
		}
		msg += v
	}

	return msg
}

func TranslateTag(fd validator.FieldError) string {
	field := fd.StructField()

	switch fd.Tag() {
	case "required":
		return fmt.Sprintf("field %s wajib diisi", field)
	case "min":
		return fmt.Sprintf("field %s size minimal %s", field, fd.Param())
	case "password":
		return fmt.Sprintf("field %s harus kombinasi harus besar, kecil, angka, dan simbol", field)
	case "eqfield":
		return fmt.Sprintf("field %s harus sama dengan %s", field, fd.Param())
	}

	return fmt.Sprintf("field %s tidak valid", field)
}
