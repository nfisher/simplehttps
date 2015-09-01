package main_test

import (
	"strings"
	"testing"

	. "."
)

const validMapping = `
{
	"apps": {
		"/path1": "http://localhost:8080",
		"/path2": "http://localhost:8081",
		"/path3": "http://localhost:8082"
	}
}
`

func Test_DecodeConfig_should_yield_error_with_invalid_json(t *testing.T) {
	r := strings.NewReader(validMapping[:len(validMapping)-3])
	c := NewConfig()

	err := DecodeConfig(r, c)

	if err == nil {
		t.Fatal("want error, got nil")
	}
}

func Test_DecodeConfig_should_yield_error_with_invalid_url(t *testing.T) {
	invalidUrlMapping := `{ "apps": { "/path1": "://localhost:8080" } }`
	r := strings.NewReader(invalidUrlMapping)
	c := NewConfig()

	err := DecodeConfig(r, c)

	if err == nil {
		t.Fatal("want error, got nil")
	}
}

func Test_DecodeConfig_should_parse_valid_json(t *testing.T) {
	r := strings.NewReader(validMapping)
	c := NewConfig()

	err := DecodeConfig(r, c)

	if err != nil {
		t.Fatalf("want err = nil, got %q", err)
	}

	if c.Len() != 3 {
		t.Fatalf("want 3, got %v", c.Len())
	}

	u := c.UrlFor("/path1")
	if u.String() != "http://localhost:8080" {
		t.Fatalf("want http://localhost:8080, got %q", u.String())
	}
}

func Test_UrlFor_should_return_nil_if_path_not_found(t *testing.T) {
	c := NewConfig()

	u := c.UrlFor("/path1")
	if u != nil {
		t.Fatalf("want u = nil, got %q", u.String())
	}
}
