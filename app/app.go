package app

	webTaskController := controllers.InitWebTaskController(taskService, catRepo, statRepo)
	webRouter := mux.NewRouter()
	routers.RegisterWebRoutes(webRouter, webTaskController)