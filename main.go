package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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

// Config
type Config struct {
	sync.RWMutex
	apps map[string]*url.URL
}

// NewConfig
func NewConfig() *Config {
	return &Config{apps: make(map[string]*url.URL)}
}

// Len
func (c *Config) Len() int {
	return len(c.apps)
}

// UrlFor
func (c *Config) UrlFor(host, path string) (u *url.URL) {
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

// DecodeConfig
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

		config.apps[p] = u
	}

	return nil
}

// Rewriter
type Rewriter struct {
	Delegate http.Handler
	Config   *Config
}

// ServeHTTP
func (rw *Rewriter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	writer := &Writer{
		0,
		w,
	}

	u := rw.Config.UrlFor(r.Host, r.URL.Path)

	if u != nil {
		// TODO: (NF 2015-09-01) Investigate whether the proxy can be cached and reused indefinitely.
		proxy := httputil.NewSingleHostReverseProxy(u)
		proxy.ServeHTTP(writer, r)
	} else {
		rw.Delegate.ServeHTTP(writer, r)
	}

	log.Printf("%v\t%v\t%q\t%v\t%v\n", r.RemoteAddr, r.Method, r.URL, time.Now().Sub(start), writer.TotalBytes())
}

// Writer
type Writer struct {
	Counter int
	http.ResponseWriter
}

// Write
func (w *Writer) Write(b []byte) (int, error) {
	count, err := w.ResponseWriter.Write(b)
	if err != nil {
		return 0, err
	}

	w.Counter += count

	return count, nil
}

// TotalBytes
func (w *Writer) TotalBytes() int {
	return w.Counter
}

// main
func main() {
	var listenAddr string
	var certPath string
	var keyPath string
	var siteRoot string
	var configFile string

	flag.StringVar(&configFile, "config", "config.json", "configuration file for application mappings.")
	flag.StringVar(&listenAddr, "listen", "127.0.0.1:8443", "listening address")
	flag.StringVar(&certPath, "cert", "certs/server.crt", "certificate path")
	flag.StringVar(&keyPath, "key", "certs/server.key", "key path")
	flag.StringVar(&siteRoot, "root", "_site", "site root directory")

	flag.Parse()

	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	c := NewConfig()
	err = DecodeConfig(file, c)
	if err != nil {
		log.Fatal(err.Error())
	}

	handler := &Rewriter{
		Delegate: http.FileServer(http.Dir(siteRoot)),
		Config:   c,
	}

	log.Printf("server listening on https://%v serving from %v\n", listenAddr, siteRoot)
	log.Fatal(http.ListenAndServeTLS(listenAddr, certPath, keyPath, handler))
}
