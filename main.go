package main

import (
	"flag"
	"log"
	"net/http"
)

const (
	defaultAddr     = "127.0.0.1:8443"
	defaultCertPath = "certs/server.crt"
	defaultKeyPath  = "certs/server.key"
	defaultSiteRoot = "_site"
)

func main() {
	var listenAddr string
	var certPath string
	var keyPath string
	var siteRoot string

	flag.StringVar(&listenAddr, "listen", defaultAddr, "listening address (default "+defaultAddr+")")
	flag.StringVar(&certPath, "cert", defaultCertPath, "certificate path (default "+defaultCertPath+")")
	flag.StringVar(&keyPath, "key", defaultKeyPath, "key path (default "+defaultCertPath+")")
	flag.StringVar(&siteRoot, "root", defaultSiteRoot, "site root directory (default "+defaultSiteRoot+")")
	flag.Parse()

	log.Printf("server listening on https://%v serving from %v\n", listenAddr, siteRoot)
	log.Fatal(http.ListenAndServeTLS(listenAddr, certPath, keyPath, http.FileServer(http.Dir(siteRoot))))
}
