package controllers

import (
	"html/template"
	"net/http"
	"rompelago/auth"
	"rompelago/middleware"
	"rompelago/models"
	"rompelago/services"
	"strconv"

	"github.com/gorilla/mux"
)

type WebMessageControllers struct {
	postService *services.PostDiscussionService
	templates   *template.Template
}

func InitWebMessageController(postService *services.PostDiscussionService) *WebMessageControllers {
	fonctionsDisponiblesDansLesTemplates := template.FuncMap{
		"additionner": func(a, b int) int { return a + b },
		"soustraire":  func(a, b int) int { return a - b },
	}
	tmpl := template.Must(template.New("").Funcs(fonctionsDisponiblesDansLesTemplates).ParseGlob("templates/*.html"))
	return &WebMessageControllers{
		postService: postService,
		templates:   tmpl,
	}
}

type UpdateMessagePageData struct {
	Post models.PostModel
}

// GET /fil/{id}/message/{msgId}/update
func (c *WebMessageControllers) UpdateMessagePage(w http.ResponseWriter, r *http.Request) {
	msgId, err := strconv.Atoi(mux.Vars(r)["msgId"])
	if err != nil || msgId <= 0 {
		http.Error(w, "Identifiant de message invalide", http.StatusBadRequest)
		return
	}

	post, err := c.postService.ReadByIdFull(msgId)
	if err != nil {
		http.Error(w, "Message introuvable : "+err.Error(), http.StatusNotFound)
		return
	}

	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		http.Redirect(w, r, "/connection", http.StatusSeeOther)
		return
	}
	if claims.UserID != strconv.Itoa(post.Creator.Id) {
		http.Error(w, "Accès interdit", http.StatusForbidden)
		return
	}

	if err := c.templates.ExecuteTemplate(w, "updatemessage", UpdateMessagePageData{Post: post}); err != nil {
		http.Error(w, "Erreur rendu template : "+err.Error(), http.StatusInternalServerError)
	}
}

// POST /fil/{id}/message/{msgId}/update
func (c *WebMessageControllers) UpdateMessageAction(w http.ResponseWriter, r *http.Request) {
	filId := mux.Vars(r)["id"]
	msgId, err := strconv.Atoi(mux.Vars(r)["msgId"])
	if err != nil || msgId <= 0 {
		http.Error(w, "Identifiant de message invalide", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		http.Redirect(w, r, "/connection", http.StatusSeeOther)
		return
	}

	post, err := c.postService.ReadByIdFull(msgId)
	if err != nil {
		http.Error(w, "Message introuvable", http.StatusNotFound)
		return
	}
	if claims.UserID != strconv.Itoa(post.Creator.Id) {
		http.Error(w, "Accès interdit", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulaire invalide", http.StatusBadRequest)
		return
	}

	if err := c.postService.Update(msgId, r.FormValue("name"), r.FormValue("contenu")); err != nil {
		http.Error(w, "Erreur lors de la modification : "+err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/fil/"+filId, http.StatusSeeOther)
}

// POST /fil/{id}/message/{msgId}/delete
func (c *WebMessageControllers) DeleteMessageAction(w http.ResponseWriter, r *http.Request) {
	filId := mux.Vars(r)["id"]
	msgId, err := strconv.Atoi(mux.Vars(r)["msgId"])
	if err != nil || msgId <= 0 {
		http.Error(w, "Identifiant de message invalide", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		http.Redirect(w, r, "/connection", http.StatusSeeOther)
		return
	}

	post, err := c.postService.ReadByIdFull(msgId)
	if err != nil {
		http.Error(w, "Message introuvable", http.StatusNotFound)
		return
	}
	if claims.UserID != strconv.Itoa(post.Creator.Id) {
		http.Error(w, "Accès interdit", http.StatusForbidden)
		return
	}

	if err := c.postService.Delete(msgId); err != nil {
		http.Error(w, "Erreur lors de la suppression : "+err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/fil/"+filId, http.StatusSeeOther)
}

func (c *WebMessageControllers) CreateMessageAction(w http.ResponseWriter, r *http.Request) {
	filId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || filId <= 0 {
		http.Error(w, "Identifiant de fil invalide", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		http.Redirect(w, r, "/connection", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulaire invalide", http.StatusBadRequest)
		return
	}

	action := r.FormValue("action")

	// --- Vote like/dislike ---
	if action == "like" || action == "dislike" {
		messageId, err := strconv.Atoi(r.FormValue("messageId"))
		if err != nil || messageId <= 0 {
			http.Error(w, "Identifiant de message invalide", http.StatusBadRequest)
			return
		}

		userId, _ := strconv.Atoi(claims.UserID) // claims déjà récupéré plus haut

		if err := c.postService.LikeDislike(userId, messageId, action); err != nil {
			http.Redirect(w, r, "/fil/"+mux.Vars(r)["id"], http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/fil/"+mux.Vars(r)["id"], http.StatusSeeOther)
		return
	}

	userId, err := strconv.Atoi(claims.UserID)
	if err != nil || userId <= 0 {
		http.Error(w, "Utilisateur invalide", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulaire invalide", http.StatusBadRequest)
		return
	}

	contenu := r.FormValue("contenu")
	if contenu == "" {
		http.Error(w, "Le contenu est obligatoire", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Le nom est obligatoire", http.StatusBadRequest)
		return
	}

	post := models.PostModel{
		Name:       name,
		Contenu:    contenu,
		Creator:    models.Utilisateur{Id: userId},
		FilAssocié: models.FilDiscussionModel{Id: filId},
	}

	if _, err := c.postService.Create(post); err != nil {
		http.Error(w, "Erreur lors de l'ajout du message : "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/fil/"+mux.Vars(r)["id"], http.StatusSeeOther)
}
