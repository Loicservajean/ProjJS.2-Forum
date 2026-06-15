package controllers

import (
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
}

func InitWebMessageController(postService *services.PostDiscussionService) *WebMessageControllers {
	return &WebMessageControllers{postService: postService}
}

// POST /fil/{id}/message
// Ajoute un message dans un fil. Requiert d'être connecté avant.
func (c *WebMessageControllers) CreateMessageAction(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'id du fil depuis l'URL de la page
	filId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || filId <= 0 {
		http.Error(w, "Identifiant de fil invalide", http.StatusBadRequest)
		return
	}

	// Récupérer l'utilisateur depuis le token JWT qui est dans le middleware.
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok || claims == nil {
		http.Redirect(w, r, "/connection", http.StatusSeeOther)
		return
	}

	userId, err := strconv.Atoi(claims.UserID)
	if err != nil || userId <= 0 {
		http.Error(w, "Utilisateur invalide", http.StatusUnauthorized)
		return
	}

	// Parser le formulaire
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

	// Rediriger vers la page du fil après soumission
	http.Redirect(w, r, "/fil/"+mux.Vars(r)["id"], http.StatusSeeOther)
}
