package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	passwordMinLength = 8
	passwordMaxLength = 32
	passwordMinLower  = 1
	passwordMinUpper  = 1
	passwordMinDigit  = 1
	passwordMinSymbol = 1
	loginMinLength    = 8
)

var (
	lengthRegexp    = regexp.MustCompile(fmt.Sprintf(`^.{%d,%d}$`, passwordMinLength, passwordMaxLength))
	lowerCaseRegexp = regexp.MustCompile(fmt.Sprintf(`[a-z]{%d,}`, passwordMinLower))
	upperCaseRegexp = regexp.MustCompile(fmt.Sprintf(`[A-Z]{%d,}`, passwordMinUpper))
	digitRegexp     = regexp.MustCompile(fmt.Sprintf(`[0-9]{%d,}`, passwordMinDigit))
	symbolRegexp    = regexp.MustCompile(fmt.Sprintf(`[!@#$%%^&*]{%d,}`, passwordMinSymbol))
	loginRegexp     = regexp.MustCompile(fmt.Sprintf(`^[a-zA-Z0-9]{%d,}$`, loginMinLength))
	digitInLogin    = regexp.MustCompile(`[0-9]`)
)

type CustomValidator struct {
	v         *validator.Validate
	passwdErr error
	loginErr  error
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()
	cv := &CustomValidator{v: v}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := v.RegisterValidation("password", cv.passwordValidate)
	if err != nil {
		panic(err)
	}

	err = v.RegisterValidation("login", cv.loginValidate)
	if err != nil {
		panic(err)
	}

	return cv
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.v.Struct(i)
	if err != nil {
		fieldErr := err.(validator.ValidationErrors)[0]

		return cv.newValidationError(fieldErr.Field(), fieldErr.Value(), fieldErr.Tag(), fieldErr.Param())
	}
	return nil
}

func (cv *CustomValidator) newValidationError(field string, value interface{}, tag string, param string) error {
	switch tag {
	case "required":
		return fmt.Errorf("field %s is required", field)
	case "password":
		return cv.passwdErr
	case "login":
		return cv.loginErr
	case "min":
		return fmt.Errorf("field %s must be at least %s characters", field, param)
	case "max":
		return fmt.Errorf("field %s must be at most %s characters", field, param)
	default:
		return fmt.Errorf("field %s is invalid", field)
	}
}

func (cv *CustomValidator) passwordValidate(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.String {
		cv.passwdErr = fmt.Errorf("field %s must be a string", fl.FieldName())
		return false
	}

	fieldValue := fl.Field().String()

	if ok := lengthRegexp.MatchString(fieldValue); !ok {
		cv.passwdErr = fmt.Errorf("field %s must be between %d and %d characters", fl.FieldName(), passwordMinLength, passwordMaxLength)
		return false
	} else if ok = lowerCaseRegexp.MatchString(fieldValue); !ok {
		cv.passwdErr = fmt.Errorf("field %s must contain at least %d lowercase letter(s)", fl.FieldName(), passwordMinLower)
		return false
	} else if ok = upperCaseRegexp.MatchString(fieldValue); !ok {
		cv.passwdErr = fmt.Errorf("field %s must contain at least %d uppercase letter(s)", fl.FieldName(), passwordMinUpper)
		return false
	} else if ok = digitRegexp.MatchString(fieldValue); !ok {
		cv.passwdErr = fmt.Errorf("field %s must contain at least %d digit(s)", fl.FieldName(), passwordMinDigit)
		return false
	} else if ok = symbolRegexp.MatchString(fieldValue); !ok {
		cv.passwdErr = fmt.Errorf("field %s must contain at least %d special character(s)", fl.FieldName(), passwordMinSymbol)
		return false
	}

	return true
}

func (cv *CustomValidator) loginValidate(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.String {
		cv.loginErr = fmt.Errorf("field %s must be a string", fl.FieldName())
		return false
	}

	fieldValue := fl.Field().String()

	if ok := loginRegexp.MatchString(fieldValue); !ok {
		cv.loginErr = fmt.Errorf("field %s must be at least %d characters long and contain only latin letters and digits", fl.FieldName(), loginMinLength)
		return false
	}

	if ok := digitInLogin.MatchString(fieldValue); !ok {
		cv.loginErr = fmt.Errorf("field %s must contain at least one digit", fl.FieldName())
		return false
	}

	return true
}
