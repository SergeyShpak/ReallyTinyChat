package main

import (
	"log"
	"net/http"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/handlers"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/routing"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Run failed with an error: %v", err)
	}
}

func run() error {
	r := routing.NewRouter()
	r.HandleFunc("/conn", routing.HttpMethodPost, handlers.Connect)
	srv := getDefaultServer(r)
	err := srv.ListenAndServe()
	return err
}

func getDefaultServer(r http.Handler) *http.Server {
	return &http.Server{
		Handler: r,
		Addr:    ":8080",
	}
}
