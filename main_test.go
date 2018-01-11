package simplehttps_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/nfisher/simplehttps"
)

const validMapping = `
{
	"apps": {
		"/path1": "http://localhost:8080",
		"/path2": "http://localhost:8081",
		"/path3/subpath1": "http://localhost:8082",
		"/path3/subpath1/subpath2": "http://localhost:8083",
		"//dev.local/path3/subpath1/subpath2": "http://localhost:8084"
	}
}
`

func Test_Write_should_increment_TotalBytes_correctly(t *testing.T) {
	w := &Writer{0, httptest.NewRecorder()}

	w.Write([]byte("hello world"))
	w.Write([]byte("!"))

	if w.TotalBytes() != 12 {
		t.Fatalf("want w.TotalBytes() = 12, got %v", w.TotalBytes())
	}
}

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

	if c.Len() < 1 {
		t.Fatalf("want c.Len() > 0, got %v", c.Len())
	}
}

func Test_UrlFor_should_return_expected_urls_for_a_given_path(t *testing.T) {
	r := strings.NewReader(validMapping)
	c := NewConfig()

	DecodeConfig(r, c)

	var testData = []struct {
		Host string
		Path string
		Url  string
	}{
		{"localhost", "/path1", "http://localhost:8080"},
		{"localhost", "/path2/index.html", "http://localhost:8081"},
		{"localhost", "/path3/subpath1/index.html", "http://localhost:8082"},
		{"dev.local", "/path3/subpath1/subpath2/index.html", "http://localhost:8084"},
	}

	for i, v := range testData {
		u := c.URLFor(v.Host, v.Path)
		if u == nil {
			t.Errorf("[%v] want u = *url.URL, got nil", i)
			continue
		}

		if u.String() != v.Url {
			t.Errorf("[%v] want u.String() = %q, got %q", i, v.Url, u.String())
		}
	}
}

func Test_UrlFor_should_return_nil_if_path_not_found(t *testing.T) {
	c := NewConfig()

	u := c.URLFor("localhost", "/path1")
	if u != nil {
		t.Fatalf("want u = nil, got %q", u.String())
	}
}
