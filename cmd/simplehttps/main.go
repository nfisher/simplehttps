package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/nfisher/simplehttps"
)

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

	c := simplehttps.NewConfig()
	err = simplehttps.DecodeConfig(file, c)
	if err != nil {
		log.Fatal(err.Error())
	}

	handler := &simplehttps.Rewriter{
		Delegate: http.FileServer(http.Dir(siteRoot)),
		Config:   c,
	}

	log.Printf("server listening on https://%v serving from %v\n", listenAddr, siteRoot)
	log.Fatal(http.ListenAndServeTLS(listenAddr, certPath, keyPath, handler))
}
