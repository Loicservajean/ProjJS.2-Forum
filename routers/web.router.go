package routers

import (
	"rompelago/controllers"

	"github.com/gorilla/mux"
)

func RegisterWebRoutes(r *mux.Router, tc *controllers.WebFilsControllers) {
	r.HandleFunc("/", nil)
}
