package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nfisher/simplehttps"
)

// Version is the injected version of the app based on the GIT SHA.
var Version = "dev"

type cmdConfig struct {
	listenAddr string
	certPath   string
	keyPath    string
	siteRoot   string
	configFile string
	version    bool
}

// main
func main() {
	var cfg cmdConfig

	flag.Usage = func() {
		fmt.Printf("Usage of simplehttps (%s):\n\n", Version)
		flag.PrintDefaults()
	}

	flag.StringVar(&cfg.configFile, "config", "config.json", "configuration file for application mappings.")
	flag.StringVar(&cfg.listenAddr, "listen", "127.0.0.1:8443", "listening address")
	flag.StringVar(&cfg.certPath, "cert", "certs/server.crt", "certificate path")
	flag.StringVar(&cfg.keyPath, "key", "certs/server.key", "key path")
	flag.StringVar(&cfg.siteRoot, "root", "_site", "site root directory")
	flag.BoolVar(&cfg.version, "version", false, "output simplehttps version")

	flag.Parse()

	if cfg.version {
		fmt.Println(Version)
		return
	}

	file, err := os.Open(cfg.configFile)
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
		Delegate: http.FileServer(http.Dir(cfg.siteRoot)),
		Config:   c,
	}

	log.Printf("server listening on https://%v serving from %v\n", cfg.listenAddr, cfg.siteRoot)
	log.Fatal(http.ListenAndServeTLS(cfg.listenAddr, cfg.certPath, cfg.keyPath, handler))
}
