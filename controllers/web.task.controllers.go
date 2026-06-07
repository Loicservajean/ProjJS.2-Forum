package controllers

import (
	"html/template"
)

type WebTaskControllers struct {
	service   *services.TaskService
	templates *template.Template
	catRepo   *repositories.CategoryRepositories
	statRepo  *repositories.StatusRepositories
}

func InitWebTaskController(service *services.TaskService, catRepo *repositories.CategoryRepositories, statRepo *repositories.StatusRepositories) *WebTaskControllers {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	return &WebTaskControllers{
		service:   service,
		templates: tmpl,
		catRepo:   catRepo,
		statRepo:  statRepo,
	}
}
