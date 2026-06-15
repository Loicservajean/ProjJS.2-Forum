package app

import (
	"database/sql"
	"log"
	"net/http"

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
	ApiRouter *mux.Router
}

func InitApp() *App {
	config.LoadEnv()

	db := config.InitDbContext()

	// Repositories
	filsRepo := repositories.InitFilsRepositories(db)
	catRepo := repositories.InitCategoryRepositories(db)
	statRepo := repositories.InitStatusRepositories(db)
	userRepo := repositories.InitUserRepositories(db)

	// Services
	filsService := services.InitFilDiscussionService(filsRepo)
	authService := services.InitAuthService(userRepo)

	// --- Serveur Web (port 8081) ---
	webFilsController := controllers.InitWebFilsController(filsService, catRepo, statRepo)
	webRouter := mux.NewRouter()
	webAuthController := controllers.InitWebAuthController(authService)
	routers.RegisterWebRoutes(webRouter, webFilsController, webAuthController)

	// --- Serveur API (port 8080) ---
	apiAuthController := controllers.InitApiAuthController(authService)
	apiFilsController := controllers.InitApiFilsController(filsService)
	apiRouter := mux.NewRouter()
	routers.RegisterApiRoutes(apiRouter, apiAuthController, apiFilsController)

	return &App{
		Db:        db,
		WebRouter: webRouter,
		ApiRouter: apiRouter,
	}
}

func (a *App) Start() {
	// Lancement du serveur API
	go func() {
		log.Printf("API lancée : http://localhost:8080")
		if err := http.ListenAndServe(":8080", a.ApiRouter); err != nil {
			log.Fatalf("Erreur lancement API - %s", err.Error())
		}
	}()

	// Serveur Web
	log.Printf("Web lancé  : http://localhost:8081")
	if err := http.ListenAndServe(":8081", a.WebRouter); err != nil {
		log.Fatalf("Erreur lancement serveur web - %s", err.Error())
	}
}

func (a *App) Close() {
	if a.Db != nil {
		a.Db.Close()
	}
}
