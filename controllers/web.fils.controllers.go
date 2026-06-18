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
	"strings"
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
	Fil           models.FilDiscussionFull
	Messages      []models.PostModel
	Connected     bool
	Page          int
	Limit         int
	Sort          string
	TotalPages    int
	ConnectedUser string
}

type UpdateFilPageData struct {
	Fil        models.FilDiscussionFull
	Categories []models.CategoriesDiscussion
	Statuts    []models.StatusModel
}

type FilListPageData struct {
	Fils          []models.FilDiscussionFull
	Page          int
	Limit         int
	TotalPages    int
	ConnectedUser string
	Categories    []models.CategoriesDiscussion
}

// Voilà les seules valeurs de "limit" qu'on accepte depuis l'URL. Si quelqu'un utilise des valeurs bizarre, on l'ignore.
var limitesAutorisees = map[int]bool{10: true, 20: true, 30: true}

func InitWebFilsController(service *services.FilDiscussionService, postService *services.PostDiscussionService, catRepo *repositories.CategoryRepositories, statRepo *repositories.StatusRepositories) *WebFilsControllers {
	fonctionsDisponiblesDansLesTemplates := template.FuncMap{
		// Fonctions de calcul pour la page précédente et la suivante.
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

// limiteEtPageDepuisRequete va voir ce que la personne a demandé dans l'URL et vérifie que ça tient la route.
// La fonction est utilisée par ListPage et DetailPage, comme ça on ne répète pas deux fois le même code.
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

// cette fonction fait le calcul pour savoir où on en est : combien de pages au total, et combien de lignes il faut sauter en SQL pour tomber sur la bonne page.
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

func sortDepuisRequete(r *http.Request) string {
	sort := r.URL.Query().Get("sort")
	if sort == "popularity" || sort == "popularite" {
		return "popularity"
	}
	return ""
}

func (c *WebFilsControllers) ListPage(w http.ResponseWriter, r *http.Request) {
	limite, page := limiteEtPageDepuisRequete(r)

	userID := ""
	connectedUserID := 0
	if claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims); ok {
		userID = claims.UserID
		if id, err := strconv.Atoi(claims.UserID); err == nil {
			connectedUserID = id
		}
	}

	total, err := c.service.CountFils(connectedUserID)
	if err != nil {
		http.Error(w, "Erreur comptage : "+err.Error(), http.StatusInternalServerError)
		return
	}

	decalage, totalPages, page := calculerPagination(limite, page, total)

	fils, err := c.service.ReadAllFull(limite, decalage, connectedUserID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des fils : "+err.Error(), http.StatusInternalServerError)
		return
	}

	categories, err := c.catRepo.ReadAll()
	if err != nil {
		http.Error(w, "Erreur chargement catégories : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Filtrage côté serveur si une catégorie est demandée
	categorieID := 0
	if catStr := r.URL.Query().Get("categorie_id"); catStr != "" {
		if v, err := strconv.Atoi(catStr); err == nil {
			categorieID = v
		}
	}
	if categorieID > 0 {
		var filtered []models.FilDiscussionFull
		for _, f := range fils {
			if f.CategorieId == categorieID {
				filtered = append(filtered, f)
			}
		}
		fils = filtered
	}

	donnees := FilListPageData{Fils: fils, Page: page, Limit: limite, TotalPages: totalPages, ConnectedUser: userID, Categories: categories}
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

	//Je récupère l'utilisateur
	creatorID := 0
	if claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims); ok {
		if id, err := strconv.Atoi(claims.UserID); err == nil {
			creatorID = id
		}
	}

	categorieID := 0
	if cat, err := strconv.Atoi(r.FormValue("categorie_id")); err == nil {
		categorieID = cat
	}

	fils := models.FilDiscussionFull{
		Name:         r.FormValue("titre"),
		Description:  r.FormValue("description"),
		DateCreation: dateCreation,
		CreatorID:    creatorID,
		CategorieId:  categorieID,
	}

	id, err := c.service.Create(fils)
	if err != nil {
		http.Error(w, "Erreur lors de la création : "+err.Error(), http.StatusBadRequest)
		return
	}

	// Statut "Ouvert" (id=1) par défaut
	if err := c.service.FilsRepository.SetStatus(id, 1); err != nil {
		http.Error(w, "Erreur statut : "+err.Error(), http.StatusInternalServerError)
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

	fil, err := c.service.ReadByIdWithDetails(id)
	if err != nil {
		http.Error(w, "Fils introuvable : "+err.Error(), http.StatusNotFound)
		return
	}
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		http.Redirect(w, r, "/connection", http.StatusSeeOther)
		return
	}
	if claims.UserID != strconv.Itoa(fil.Creator.Id) {
		http.Error(w, "Accès interdit", http.StatusForbidden)
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

	sort := sortDepuisRequete(r)
	messages, err := c.postService.ReadByFilId(id, limite, decalage, sort)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des messages : "+err.Error(), http.StatusInternalServerError)
		return
	}

	userID := ""
	claims, connected := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if connected {
		userID = claims.UserID
	}

	data := FilDetailPageData{
		Fil:           fils,
		Messages:      messages,
		Connected:     connected,
		Page:          page,
		Limit:         limite,
		Sort:          sort,
		TotalPages:    totalPages,
		ConnectedUser: userID,
	}

	if err := c.templates.ExecuteTemplate(w, "fils", data); err != nil {
		http.Error(w, "Erreur rendu template : "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *WebFilsControllers) UpdateFil(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Identifiant invalide", http.StatusBadRequest)
		return
	}
	fil, err := c.service.ReadByIdWithDetails(id)
	if err != nil {
		http.Error(w, "Fils introuvable : "+err.Error(), http.StatusNotFound)
		return
	}
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
	data := UpdateFilPageData{Fil: fil, Categories: categories, Statuts: statuts}
	if err := c.templates.ExecuteTemplate(w, "updatefil", data); err != nil {
		http.Error(w, "Erreur rendu template : "+err.Error(), http.StatusInternalServerError)
	}
}

func (c *WebFilsControllers) UpdateAction(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Identifiant invalide", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulaire invalide", http.StatusBadRequest)
		return
	}

	dateCreation := r.FormValue("date_creation")
	if dateCreation == "" {
		dateCreation = time.Now().Format("2006-01-02")
	}

	categorieID := 0
	if cat, err := strconv.Atoi(r.FormValue("categorie_id")); err == nil {
		categorieID = cat
	}

	statusID := 0
	if st, err := strconv.Atoi(r.FormValue("status_id")); err == nil {
		statusID = st
	}

	fils := models.FilDiscussionFull{
		Description:  r.FormValue("description"),
		DateCreation: dateCreation,
		CategorieId:  categorieID,
		TagId:        statusID,
	}

	if err := c.service.UpdateFull(id, fils); err != nil {
		http.Error(w, "Erreur lors de la mise à jour : "+err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/forum", http.StatusSeeOther)
}

type SearchResultPageData struct {
	Fils          []models.FilDiscussionFull
	MotCle        string
	Page          int
	Limit         int
	TotalPages    int
	ConnectedUser string
}

// SearchPage vérifi la recherche d'un fil de discussion, la route elle est protégée par RequireAuth donc seul un utilisateur connecté arrive ici
func (c *WebFilsControllers) SearchPage(w http.ResponseWriter, r *http.Request) {
	// On récupère ce que l'utilisateur a tapé dans la barre de recherche
	motCle := strings.TrimSpace(r.URL.Query().Get("q"))

	userID := ""
	if claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims); ok {
		userID = claims.UserID
	}

	// Si la recherche est vide,la page est afficher sans résultat plutôt que de planter
	if motCle == "" {
		data := SearchResultPageData{MotCle: "", Page: 1, Limit: 10, TotalPages: 1, ConnectedUser: userID}
		if err := c.templates.ExecuteTemplate(w, "recherche", data); err != nil {
			http.Error(w, "Erreur rendu template : "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	limite, page := limiteEtPageDepuisRequete(r)

	total, err := c.service.CountSearch(motCle)
	if err != nil {
		http.Error(w, "Erreur comptage recherche : "+err.Error(), http.StatusInternalServerError)
		return
	}

	decalage, totalPages, page := calculerPagination(limite, page, total)

	fils, err := c.service.Search(motCle, limite, decalage)
	if err != nil {
		http.Error(w, "Erreur lors de la recherche : "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := SearchResultPageData{
		Fils:          fils,
		MotCle:        motCle,
		Page:          page,
		Limit:         limite,
		TotalPages:    totalPages,
		ConnectedUser: userID,
	}
	if err := c.templates.ExecuteTemplate(w, "recherche", data); err != nil {
		http.Error(w, "Erreur rendu template : "+err.Error(), http.StatusInternalServerError)
	}
}