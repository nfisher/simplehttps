package main_test

import (
	"encoding/json"
	"io"
	"net/url"
	"strings"
	"testing"
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

type Config struct {
	Apps map[string]*url.URL
}

func (c *Config) Len() int {
	return len(c.Apps)
}

func DecodeConfig(r io.Reader) (config *Config, err error) {
	dec := json.NewDecoder(r)
	config = &Config{}

	err = dec.Decode(config)

	return config, err
}

func Test_should_parse_valid_json(t *testing.T) {
	r := strings.NewReader(validMapping)

	c, _ := DecodeConfig(r)

	if c == nil {
		t.Fatal("want a *Config, got nil")
	}

	if c.Len() != 3 {
		t.Fatalf("want 3, got %q", c.Len())
	}

	u, ok := c.Apps["/path1"]
	if !ok {
		t.Fatal("want ok, got false")
	}

	if u == nil {
		t.Fatalf("want *URL, got nil")
	}

	if u.Host != "localhost:8080" {
		t.Fatalf("want localhost:8080, got %q", u.String())
	}
}
