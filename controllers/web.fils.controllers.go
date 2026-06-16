package controllers

import (
	"html/template"
	"net/http"
	"rompelago/auth"
	"rompelago/middleware"
	"rompelago/models"
	"rompelago/repositories"
	"rompelago/services"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type WebFilsControllers struct {
	service     *services.FilDiscussionService
	postService *services.PostDiscussionService
	templates   *template.Template
	catRepo     *repositories.CategoryRepositories
	statRepo    *repositories.StatusRepositories
}

type CreatePageData struct {
	Categories []models.CategoriesDiscussion
	Statuts    []models.StatusModel
}

type FilDetailPageData struct {
	Fil       models.FilDiscussionFull
	Messages  []models.PostModel
	Connected bool
}

func InitWebFilsController(service *services.FilDiscussionService, postService *services.PostDiscussionService, catRepo *repositories.CategoryRepositories, statRepo *repositories.StatusRepositories) *WebFilsControllers {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	return &WebFilsControllers{
		service:     service,
		postService: postService,
		templates:   tmpl,
		catRepo:     catRepo,
		statRepo:    statRepo,
	}
}

func (c *WebFilsControllers) ListPage(w http.ResponseWriter, r *http.Request) {
	fils, err := c.service.ReadAllFull()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des fils : "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := c.templates.ExecuteTemplate(w, "fils.list", fils); err != nil {
		http.Error(w, "Erreur rendu template : "+err.Error(), http.StatusInternalServerError)
	}
}

func (c *WebFilsControllers) CreateFil(w http.ResponseWriter, r *http.Request) {
	categories, err := c.catRepo.ReadAll()
	if err != nil {
		http.Error(w, "Erreur chargement catégories : "+err.Error(), http.StatusInternalServerError)
		return
	}
	statuts, err := c.statRepo.ReadAll()
	if err != nil {
		http.Error(w, "Erreur chargement statuts : "+err.Error(), http.StatusInternalServerError)
		return
	}
	data := CreatePageData{Categories: categories, Statuts: statuts}
	if err := c.templates.ExecuteTemplate(w, "nouveau", data); err != nil {
		http.Error(w, "Erreur rendu template : "+err.Error(), http.StatusInternalServerError)
	}
}

func (c *WebFilsControllers) CreateAction(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulaire invalide", http.StatusBadRequest)
		return
	}

	dateCreation := r.FormValue("date_creation")
	if dateCreation == "" {
		dateCreation = time.Now().Format("2006-01-02")
	}

	fils := models.FilDiscussionFull{
		Name:         r.FormValue("titre"),
		Description:  r.FormValue("description"),
		DateCreation: dateCreation,
	}

	if _, err := c.service.Create(fils); err != nil {
		http.Error(w, "Erreur lors de la création : "+err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/forum", http.StatusSeeOther)
}

func (c *WebFilsControllers) DeletePage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Identifiant invalide", http.StatusBadRequest)
		return
	}
	fils, err := c.service.ReadByIdWithDetails(id)
	if err != nil {
		http.Error(w, "Fils introuvable : "+err.Error(), http.StatusNotFound)
		return
	}
	if err := c.templates.ExecuteTemplate(w, "fils.delete", fils); err != nil {
		http.Error(w, "Erreur rendu template : "+err.Error(), http.StatusInternalServerError)
	}
}

func (c *WebFilsControllers) DeleteAction(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Identifiant invalide", http.StatusBadRequest)
		return
	}
	if err := c.service.Delete(id); err != nil {
		http.Error(w, "Erreur lors de la suppression : "+err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/forum", http.StatusSeeOther)
}

func (c *WebFilsControllers) DetailPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Identifiant invalide", http.StatusBadRequest)
		return
	}
	fils, err := c.service.ReadByIdWithDetails(id)
	if err != nil {
		http.Error(w, "Fils introuvable : "+err.Error(), http.StatusNotFound)
		return
	}

	messages, err := c.postService.ReadByFilId(id)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des messages : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Permet d'afficher (ou non) le formulaire de réponse selon que l'utilisateur
	// est connecté (cookie JWT valide lu par le WebAuthMiddleware).
	_, connected := r.Context().Value(middleware.UserContextKey).(*auth.Claims)

	data := FilDetailPageData{
		Fil:       fils,
		Messages:  messages,
		Connected: connected,
	}

	if err := c.templates.ExecuteTemplate(w, "fils", data); err != nil {
		http.Error(w, "Erreur rendu template : "+err.Error(), http.StatusInternalServerError)
	}
}
