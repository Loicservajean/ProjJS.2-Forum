package routers

import (
	"net/http"
	"rompelago/controllers"

	"github.com/gorilla/mux"
)

func RegisterWebRoutes(r *mux.Router, tc *controllers.WebFilsControllers) {
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	r.HandleFunc("/forum", tc.ListPage).Methods("GET")
}
