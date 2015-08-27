package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

/**
  Target JSON
	{
		"apps": {
			"/path1": ["http://localhost:8080"],
			"/path2": ["http://localhost:8081"],
			"/path3": ["http://localhost:8082"]
		}
	}
*/

// Rewriter
type Rewriter struct {
	Delegate     http.Handler
	HostMappings map[string]*url.URL
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

// ServeHTTP
func (rw *Rewriter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	writer := &Writer{
		0,
		w,
	}

	if r.URL.Path == "/" {
		googleUrl, err := url.Parse("https://www.google.co.uk/")
		if err != nil {
			http.Error(writer, "Internal Error", http.StatusInternalServerError)
			return
		}

		google := httputil.NewSingleHostReverseProxy(googleUrl)
		google.ServeHTTP(writer, r)
	} else {
		rw.Delegate.ServeHTTP(writer, r)
	}

	log.Printf("%v\t%v\t%v\t%v\t%v\n", r.RemoteAddr, r.Method, r.URL, time.Now().Sub(start), writer.TotalBytes())
}

// main
func main() {
	var listenAddr string
	var certPath string
	var keyPath string
	var siteRoot string

	flag.StringVar(&listenAddr, "listen", "127.0.0.1:8443", "listening address")
	flag.StringVar(&certPath, "cert", "certs/server.crt", "certificate path")
	flag.StringVar(&keyPath, "key", "certs/server.key", "key path")
	flag.StringVar(&siteRoot, "root", "_site", "site root directory")

	flag.Parse()

	handler := &Rewriter{
		Delegate: http.FileServer(http.Dir(siteRoot)),
	}

	log.Printf("server listening on https://%v serving from %v\n", listenAddr, siteRoot)
	log.Fatal(http.ListenAndServeTLS(listenAddr, certPath, keyPath, handler))
}
