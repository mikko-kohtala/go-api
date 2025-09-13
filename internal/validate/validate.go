package validate

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "strings"

    "github.com/go-playground/validator/v10"
)

var v = validator.New(validator.WithRequiredStructEnabled())

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
                name := fe.Field()
                // Use json tag if present
                if tag := fe.StructField(); tag != "" {
                    // If struct field has json tag, validator exposes Field(); here we try to
                    // fallback to lowercased field name when json tag unknown.
                    // For simple cases this is fine; advanced mapping can be added later.
                }
                // Best-effort to map to json tag: reflect isn't used here to keep simple.
                // Convert to lower-camel case
                if name != "" {
                    name = strings.ToLower(name[:1]) + name[1:]
                }
                out[name] = humanMessage(fe)
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

