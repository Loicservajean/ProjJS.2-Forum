package main

import (
	"log"
	"net/http"
	"rompelago/app"
)

func main() {
	application := app.InitApp()
	defer application.Close()

	log.Printf("Serveur lancé : http://localhost:8081")
	if err := http.ListenAndServe(":8081", application.WebRouter); err != nil {
		log.Fatalf("Erreur lancement serveur - %s", err.Error())
	}
}
