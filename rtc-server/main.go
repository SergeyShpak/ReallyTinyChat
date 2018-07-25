package main

import (
	"log"
	"net/http"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/cors"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/handlers"
)

func main() {
	run()
}

func run() {
	//srv := getDefaultServer(nil)
	//http.Handle("/", )
	http.HandleFunc("/conn", handlers.Connect)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getDefaultServer(r http.Handler) *http.Server {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowedHeaders: []string{
			"content-type",
			"cache-control",
		},
		ExposedHeaders: []string{
			"content-type",
		},
	})
	return &http.Server{
		Handler: corsHandler.Handler(r),
		Addr:    ":8080",
	}
}
