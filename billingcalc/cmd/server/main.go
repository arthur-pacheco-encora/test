package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/api/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	r := httprouter.New()
	routes.Install(r)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
