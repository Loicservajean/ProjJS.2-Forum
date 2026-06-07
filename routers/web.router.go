package routers

import "github.com/gorilla/mux"

func RegisterWebRoutes(r *mux.Router) {
	r.HandleFunc("/", nil)
}
