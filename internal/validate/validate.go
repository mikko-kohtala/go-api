package validate

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "reflect"
    "strings"

    "github.com/go-playground/validator/v10"
)

var v = validator.New(validator.WithRequiredStructEnabled())

func init() {
    // Use JSON tag names in validation errors
    v.RegisterTagNameFunc(func(fld reflect.StructField) string {
        name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
        if name == "-" || name == "" {
            return fld.Name
        }
        return name
    })
}

// Errors represents field validation errors keyed by JSON field name.
type Errors map[string]string

// BindAndValidate decodes JSON into dst (disallowing unknown fields) and validates it.
func BindAndValidate(r *http.Request, dst any) (Errors, error) {
    if r.Body == nil {
        return nil, errors.New("empty body")
    }
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()
    if err := dec.Decode(dst); err != nil {
        return nil, err
    }
    if err := v.Struct(dst); err != nil {
        if verrs, ok := err.(validator.ValidationErrors); ok {
            out := Errors{}
            for _, fe := range verrs {
                out[fe.Field()] = humanMessage(fe)
            }
            return out, nil
        }
        return nil, err
    }
    return nil, nil
}

func humanMessage(fe validator.FieldError) string {
    switch fe.Tag() {
    case "required":
        return "is required"
    case "min":
        return fmt.Sprintf("must be at least %s", fe.Param())
    case "max":
        return fmt.Sprintf("must be at most %s", fe.Param())
    case "oneof":
        return fmt.Sprintf("must be one of %s", fe.Param())
    case "email":
        return "must be a valid email"
    default:
        return fmt.Sprintf("is invalid (%s)", fe.Tag())
    }
}
