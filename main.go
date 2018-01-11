package simplehttps

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

/*
  Target JSON
	{
		"apps": {
			"/path1": ["http://localhost:8080"],
			"/path2": ["http://localhost:8081"],
			"/path3": ["http://localhost:8082"]
		}
	}
*/

// rawConfig is an intermediate transform for the raw JSON configuration.
type rawConfig struct {
	Apps map[string]string
}

// Config maps API root paths to URL backends.
type Config struct {
	sync.RWMutex
	apps map[string]*url.URL
}

// NewConfig generates a runtime configuration of the path to backend URL
// mapping.
func NewConfig() *Config {
	return &Config{apps: make(map[string]*url.URL)}
}

// Len returns the number of registered backend services.
func (c *Config) Len() int {
	return len(c.apps)
}

// URLFor maps the host and path to a backend URL.
func (c *Config) URLFor(host, path string) (u *url.URL) {
	u = c.urlFor("//"+host, path)
	if u != nil {
		return u
	}

	u = c.urlFor("", path)

	return u
}

func (c *Config) urlFor(host, path string) (u *url.URL) {
	baseComponents := strings.Split(path, "/")
	if host != "" {
		baseComponents = append([]string{host}, baseComponents[1:len(baseComponents)-1]...)
	}

	for i := 0; i < len(baseComponents); i++ {
		currentPath := strings.Join(baseComponents[:len(baseComponents)-i], "/")

		c.RLock()
		u, ok := c.apps[currentPath]
		c.RUnlock()

		if ok {
			return u
		}
	}

	return nil
}

// Add adds a path p and url u to the configuration.
func (c *Config) Add(p string, u *url.URL) {
	c.Lock()
	c.apps[p] = u
	c.Unlock()
}

// DecodeConfig reads the config from the reader r.
func DecodeConfig(r io.Reader, config *Config) (err error) {
	dec := json.NewDecoder(r)
	raw := &rawConfig{}

	err = dec.Decode(raw)
	if err != nil {
		return err
	}

	for p, urlString := range raw.Apps {
		u, err := url.Parse(urlString)
		if err != nil {
			return err
		}

		config.Add(p, u)
	}

	return nil
}

// Rewriter is a URL rewriter middleware that maps requests to backend services.
type Rewriter struct {
	Delegate http.Handler
	Config   *Config
}

// ServeHTTP wraps the request wih rewrite goodness and logging.
func (rw *Rewriter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	writer := &Writer{
		0,
		w,
	}

	u := rw.Config.URLFor(r.Host, r.URL.Path)

	if u != nil {
		// TODO: (NF 2015-09-01) Investigate whether the proxy can be cached and
		// reused indefinitely.
		proxy := httputil.NewSingleHostReverseProxy(u)
		r.Header.Set("x-forwarded-proto", "https")
		proxy.ServeHTTP(writer, r)
	} else {
		rw.Delegate.ServeHTTP(writer, r)
	}

	log.Printf("%v\t%v\t%q\t%v\t%v\n", r.RemoteAddr, r.Method, r.URL, time.Now().Sub(start), writer.TotalBytes())
}

// Writer counts the number of bytes written to the response.
type Writer struct {
	Counter int
	http.ResponseWriter
}

// Write composes the ResponseWriter function with a byte counter.
func (w *Writer) Write(b []byte) (int, error) {
	count, err := w.ResponseWriter.Write(b)
	if err != nil {
		return 0, err
	}

	w.Counter += count

	return count, nil
}

// TotalBytes returns the number of bytes written to the current request stream.
func (w *Writer) TotalBytes() int {
	return w.Counter
}
