package routers

import (
	"net/http"
	"rompelago/controllers"
	"rompelago/middleware"

	"github.com/gorilla/mux"
)

func RegisterWebRoutes(r *mux.Router, tc *controllers.WebFilsControllers, ac *controllers.WebAuthControllers) {
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	// Middleware non bloquant : lit le cookie JWT et injecte les claims
	// dans le contexte si présent, sans empêcher l'accès aux pages publiques.
	r.Use(middleware.WebAuthMiddleware)

	r.HandleFunc("/forum", tc.ListPage).Methods("GET")
	r.HandleFunc("/nouveau", tc.CreateFil).Methods("GET")
	r.HandleFunc("/inscription", ac.RegisterPage).Methods("GET")
	r.HandleFunc("/inscription", ac.RegisterAction).Methods("POST")
	r.HandleFunc("/connection", ac.LoginPage).Methods("GET")
	r.HandleFunc("/connection", ac.LoginAction).Methods("POST")
	r.HandleFunc("/logout", ac.LogoutAction).Methods("POST")
	r.HandleFunc("/user", ac.UserPage).Methods("GET")
}
