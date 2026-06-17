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
	Fil        models.FilDiscussionFull
	Messages   []models.PostModel
	Connected  bool
	Page       int
	Limit      int
	TotalPages int
}

type FilListPageData struct {
	Fils       []models.FilDiscussionFull
	Page       int
	Limit      int
	TotalPages int
}

// limitesAutorisees liste les seules valeurs de "limit" qu'on accepte depuis l'URL.
var limitesAutorisees = map[int]bool{10: true, 20: true, 30: true}

func InitWebFilsController(service *services.FilDiscussionService, postService *services.PostDiscussionService, catRepo *repositories.CategoryRepositories, statRepo *repositories.StatusRepositories) *WebFilsControllers {
	// additionner/soustraire sont utilisées dans les templates HTML pour calculer
	// la page précédente/suivante, puisque les templates ne savent pas faire de calcul.
	fonctionsDisponiblesDansLesTemplates := template.FuncMap{
		"additionner": func(a, b int) int { return a + b },
		"soustraire":  func(a, b int) int { return a - b },
	}
	tmpl := template.Must(template.New("").Funcs(fonctionsDisponiblesDansLesTemplates).ParseGlob("templates/*.html"))
	return &WebFilsControllers{
		service:     service,
		postService: postService,
		templates:   tmpl,
		catRepo:     catRepo,
		statRepo:    statRepo,
	}
}

// limiteEtPageDepuisRequete lit et valide les paramètres "limit" et "page" de l'URL.
// Réutilisée par ListPage et DetailPage pour éviter de dupliquer cette logique.
func limiteEtPageDepuisRequete(r *http.Request) (limite int, page int) {
	limite = 10
	valeurLimite := r.URL.Query().Get("limit")
	if valeurLimite == "tout" {
		limite = 0
	} else if limiteConvertie, err := strconv.Atoi(valeurLimite); err == nil && limitesAutorisees[limiteConvertie] {
		limite = limiteConvertie
	}

	page = 1
	if pageConvertie, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && pageConvertie > 0 {
		page = pageConvertie
	}
	return limite, page
}

// calculerPagination calcule le décalage SQL et le nombre total de pages.
// Réutilisée par ListPage (fils) et DetailPage (messages).
func calculerPagination(limite, page, total int) (decalage int, totalPages int, pageCorrigee int) {
	if limite <= 0 {
		return 0, 1, page
	}
	totalPages = (total + limite - 1) / limite
	if totalPages < 1 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}
	decalage = (page - 1) * limite
	return decalage, totalPages, page
}

func (c *WebFilsControllers) ListPage(w http.ResponseWriter, r *http.Request) {
	limite, page := limiteEtPageDepuisRequete(r)

	total, err := c.service.CountFils()
	if err != nil {
		http.Error(w, "Erreur comptage : "+err.Error(), http.StatusInternalServerError)
		return
	}

	decalage, totalPages, page := calculerPagination(limite, page, total)

	fils, err := c.service.ReadAllFull(limite, decalage)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des fils : "+err.Error(), http.StatusInternalServerError)
		return
	}

	donnees := FilListPageData{Fils: fils, Page: page, Limit: limite, TotalPages: totalPages}
	if err := c.templates.ExecuteTemplate(w, "fils.list", donnees); err != nil {
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

	limite, page := limiteEtPageDepuisRequete(r)

	total, err := c.postService.CountByFilId(id)
	if err != nil {
		http.Error(w, "Erreur comptage des messages : "+err.Error(), http.StatusInternalServerError)
		return
	}

	decalage, totalPages, page := calculerPagination(limite, page, total)

	messages, err := c.postService.ReadByFilId(id, limite, decalage)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des messages : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Permet d'afficher (ou non) le formulaire de réponse selon que l'utilisateur
	// est connecté (cookie JWT valide lu par le WebAuthMiddleware).
	_, connected := r.Context().Value(middleware.UserContextKey).(*auth.Claims)

	data := FilDetailPageData{
		Fil:        fils,
		Messages:   messages,
		Connected:  connected,
		Page:       page,
		Limit:      limite,
		TotalPages: totalPages,
	}

	if err := c.templates.ExecuteTemplate(w, "fils", data); err != nil {
		http.Error(w, "Erreur rendu template : "+err.Error(), http.StatusInternalServerError)
	}
}
