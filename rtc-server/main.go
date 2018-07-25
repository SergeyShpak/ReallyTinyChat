package main

import (
	"log"
	"net/http"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/handlers"
)

func main() {
	run()
}

func run() {
	http.HandleFunc("/conn", handlers.Connect)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
