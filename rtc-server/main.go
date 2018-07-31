package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/handlers"
)

func main() {
	run()
}

func run() {
	log.Println("Starting server")
	srv := setupServ(createHttpServer())
	var wg sync.WaitGroup
	var err error
	wg.Add(1)
	go func() {
		err = listenAndServeTLS(srv)
		wg.Done()
	}()
	log.Println("Serving")
	wg.Wait()
	if err != nil {
		log.Println("an error occurred while serving TLS: ", err)
	}
	log.Println("Stopping server")
}

func createHttpServer() *http.Server {
	r := newRouter()
	r.Add("/conn", handlers.Connect)
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      r,
		Addr:         ":4443",
	}
	return srv
}

func setupServ(srv *http.Server) *http.Server {
	srv.TLSConfig = &tls.Config{
		// Causes servers to use Go's default ciphersuite preferences,
		// which are tuned to avoid attacks. Does nothing on clients.
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
	return srv
}

func listenAndServeTLS(srv *http.Server) error {
	return srv.ListenAndServeTLS("server.crt", "server.key")
}
