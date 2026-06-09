package app

import (
	"database/sql"

	"rompelago/config"
	"rompelago/controllers"
	"rompelago/repositories"
	"rompelago/routers"
	"rompelago/services"

	"github.com/gorilla/mux"
)

type App struct {
	Db        *sql.DB
	WebRouter *mux.Router
}

func InitApp() *App {
	config.LoadEnv()

	db := config.InitDbContext()

	filsRepo := repositories.InitFilsRepositories(db)
	catRepo := repositories.InitCategoryRepositories(db)
	statRepo := repositories.InitStatusRepositories(db)

	filsService := services.InitFilDiscussionService(filsRepo)

	// Web — port 8081
	webFilsController := controllers.InitWebFilsController(filsService, catRepo, statRepo)
	webRouter := mux.NewRouter()
	routers.RegisterWebRoutes(webRouter, webFilsController)

	return &App{
		Db:        db,
		WebRouter: webRouter,
	}
}

func (a *App) Close() {
	if a.Db != nil {
		a.Db.Close()
	}
}
