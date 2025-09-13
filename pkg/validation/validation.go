package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	validate.RegisterValidation("phone", validatePhone)
	validate.RegisterValidation("username", validateUsername)
}

func Validate(s interface{}) error {
	return validate.Struct(s)
}

func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()
			param := e.Param()

			switch tag {
			case "required":
				errors[field] = fmt.Sprintf("%s is required", field)
			case "email":
				errors[field] = fmt.Sprintf("%s must be a valid email address", field)
			case "min":
				if e.Type().Kind() == reflect.String {
					errors[field] = fmt.Sprintf("%s must be at least %s characters long", field, param)
				} else {
					errors[field] = fmt.Sprintf("%s must be at least %s", field, param)
				}
			case "max":
				if e.Type().Kind() == reflect.String {
					errors[field] = fmt.Sprintf("%s must be at most %s characters long", field, param)
				} else {
					errors[field] = fmt.Sprintf("%s must be at most %s", field, param)
				}
			case "len":
				errors[field] = fmt.Sprintf("%s must be exactly %s characters long", field, param)
			case "url":
				errors[field] = fmt.Sprintf("%s must be a valid URL", field)
			case "uuid":
				errors[field] = fmt.Sprintf("%s must be a valid UUID", field)
			case "phone":
				errors[field] = fmt.Sprintf("%s must be a valid phone number", field)
			case "username":
				errors[field] = fmt.Sprintf("%s must contain only letters, numbers, and underscores", field)
			case "oneof":
				errors[field] = fmt.Sprintf("%s must be one of: %s", field, param)
			case "gt":
				errors[field] = fmt.Sprintf("%s must be greater than %s", field, param)
			case "gte":
				errors[field] = fmt.Sprintf("%s must be greater than or equal to %s", field, param)
			case "lt":
				errors[field] = fmt.Sprintf("%s must be less than %s", field, param)
			case "lte":
				errors[field] = fmt.Sprintf("%s must be less than or equal to %s", field, param)
			default:
				errors[field] = fmt.Sprintf("%s failed %s validation", field, tag)
			}
		}
	}

	return errors
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")
	phone = strings.ReplaceAll(phone, "+", "")

	if len(phone) < 10 || len(phone) > 15 {
		return false
	}

	for _, ch := range phone {
		if ch < '0' || ch > '9' {
			return false
		}
	}

	return true
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	if len(username) < 3 || len(username) > 30 {
		return false
	}

	for _, ch := range username {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_') {
			return false
		}
	}

	return true
}