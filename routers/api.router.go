package routers

import (
	"rompelago/controllers"

	"github.com/gorilla/mux"
)

func RegisterApiRoutes(r *mux.Router, authCtrl *controllers.ApiAuthController, filsCtrl *controllers.ApiFilsController) {
	api := r.PathPrefix("/api").Subrouter()

	// --- Auth ---
	api.HandleFunc("/auth/login", authCtrl.Login).Methods("POST")
	api.HandleFunc("/auth/register", authCtrl.Register).Methods("POST")

	// --- Fils ---
	api.HandleFunc("/fils", filsCtrl.ListFils).Methods("GET")
}
