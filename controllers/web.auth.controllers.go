package controllers

import (
	"html/template"
	"net/http"
	"rompelago/auth"
	"rompelago/dto"
	"rompelago/middleware"
	"rompelago/services"
	"strconv"
)

type WebAuthControllers struct {
	authService *services.AuthService
	templates   *template.Template
}

func InitWebAuthController(authService *services.AuthService) *WebAuthControllers {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	return &WebAuthControllers{
		authService: authService,
		templates:   tmpl,
	}
}

// GET /inscription
func (c *WebAuthControllers) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if err := c.templates.ExecuteTemplate(w, "inscription", map[string]string{}); err != nil {
		http.Error(w, "Erreur rendu template : "+err.Error(), http.StatusInternalServerError)
	}
}

// POST /inscription
func (c *WebAuthControllers) RegisterAction(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulaire invalide", http.StatusBadRequest)
		return
	}

	req := dto.RegisterRequestDto{
		Pseudo:   r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if r.FormValue("confirm_password") != req.Password {
		w.WriteHeader(http.StatusBadRequest)
		c.templates.ExecuteTemplate(w, "inscription", map[string]string{
			"Error": "Les mots de passe ne correspondent pas",
		})
		return
	}

	if err := c.authService.Register(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.templates.ExecuteTemplate(w, "inscription", map[string]string{
			"Error": err.Error(),
		})
		return
	}

	http.Redirect(w, r, "/connection", http.StatusSeeOther)
}

// GET /connexion
func (c *WebAuthControllers) LoginPage(w http.ResponseWriter, r *http.Request) {
	if err := c.templates.ExecuteTemplate(w, "connection", map[string]string{}); err != nil {
		http.Error(w, "Erreur rendu template : "+err.Error(), http.StatusInternalServerError)
	}
}

// POST /connexion
func (c *WebAuthControllers) LoginAction(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulaire invalide", http.StatusBadRequest)
		return
	}

	req := dto.LoginRequestDto{
		Pseudo:   r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	resp, err := c.authService.Login(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		c.templates.ExecuteTemplate(w, "connection", map[string]string{
			"Error": err.Error(),
		})
		return
	}

	// Le JWT est stocké dans un cookie HttpOnly afin d'être renvoyé
	// automatiquement par le navigateur sur les pages web suivantes.
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    resp.AccessToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   resp.ExpiresIn,
	})

	http.Redirect(w, r, "/forum", http.StatusSeeOther)
}

// POST /logout
func (c *WebAuthControllers) LogoutAction(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/connection", http.StatusSeeOther)
}

// GET /user
func (c *WebAuthControllers) UserPage(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		// Non connecté : on affiche la page avec un utilisateur vide (nil),
		// le template invite alors à se connecter.
		if err := c.templates.ExecuteTemplate(w, "user.profile", nil); err != nil {
			http.Error(w, "Erreur rendu template : "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	userID, err := strconv.Atoi(claims.UserID)
	if err != nil {
		http.Error(w, "Identifiant utilisateur invalide dans le token", http.StatusInternalServerError)
		return
	}

	user, err := c.authService.GetById(userID)
	if err != nil {
		http.Error(w, "Erreur récupération profil : "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := c.templates.ExecuteTemplate(w, "user.profile", user); err != nil {
		http.Error(w, "Erreur rendu template : "+err.Error(), http.StatusInternalServerError)
	}
}
