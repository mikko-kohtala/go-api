package validate

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type sample struct {
	Email string `json:"email" validate:"required,email"`
}

func TestBindAndValidate_UnknownField(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"email":"a@b.com","oops":1}`))
	errs, err := BindAndValidate(r, &sample{})
	if err == nil {
		t.Fatalf("expected error for unknown field, got nil and errs=%v", errs)
	}
}

func TestBindAndValidate_FieldNamesFromJSONTags(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"email":"not-an-email"}`))
	errs, err := BindAndValidate(r, &sample{})
	if err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}
	if errs == nil || errs["email"] == "" {
		t.Fatalf("expected field error keyed by 'email', got: %v", errs)
	}
}
