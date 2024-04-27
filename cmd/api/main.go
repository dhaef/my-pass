package main

import (
	"log"
	"net/http"

	"github.com/dhaef/my-pass/internal/db"
)

func main() {
	db.Connect()
	db.SetupTablesAndUser()
	s := http.Server{
		Addr:    ":3000",
		Handler: buildRoutes(),
	}

	log.Fatal(s.ListenAndServe())
}
