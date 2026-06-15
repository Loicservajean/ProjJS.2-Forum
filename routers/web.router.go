package routers

import (
	"net/http"
	"rompelago/controllers"

	"github.com/gorilla/mux"
)

func RegisterWebRoutes(r *mux.Router, tc *controllers.WebFilsControllers, ac *controllers.WebAuthControllers) {
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	r.HandleFunc("/forum", tc.ListPage).Methods("GET")
	r.HandleFunc("/inscription", ac.RegisterPage).Methods("GET")
	r.HandleFunc("/inscription", ac.RegisterAction).Methods("POST")
	r.HandleFunc("/connection", ac.LoginPage).Methods("GET")
	r.HandleFunc("/connection", ac.LoginAction).Methods("POST")
}
