package routers

import (
	"net/http"
	"rompelago/controllers"
	"rompelago/middleware"

	"github.com/gorilla/mux"
)

func RegisterWebRoutes(r *mux.Router, tc *controllers.WebFilsControllers, ac *controllers.WebAuthControllers, mc *controllers.WebMessageControllers) {
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	r.Use(middleware.WebAuthMiddleware)

	// Auth
	r.HandleFunc("/inscription", ac.RegisterPage).Methods("GET")
	r.HandleFunc("/inscription", ac.RegisterAction).Methods("POST")
	r.HandleFunc("/connection", ac.LoginPage).Methods("GET")
	r.HandleFunc("/connection", ac.LoginAction).Methods("POST")

	// Fils de discussion
	r.HandleFunc("/forum", tc.ListPage).Methods("GET")
	r.Handle("/nouveau", middleware.RequireAuth(http.HandlerFunc(tc.CreateFil))).Methods("GET")
	r.Handle("/nouveaufil", middleware.RequireAuth(http.HandlerFunc(tc.CreateAction))).Methods("POST")
	r.HandleFunc("/fil/{id}", tc.DetailPage).Methods("GET")
	r.Handle("/fil/{id}/update", middleware.RequireAuth(http.HandlerFunc(tc.UpdateFil))).Methods("GET")
	r.Handle("/fil/{id}/update", middleware.RequireAuth(http.HandlerFunc(tc.UpdateAction))).Methods("POST")

	// Messages (Normalement non accessible sans être connecté)
	r.Handle("/fil/{id}/message", middleware.RequireAuth(http.HandlerFunc(mc.CreateMessageAction))).Methods("POST")
}
