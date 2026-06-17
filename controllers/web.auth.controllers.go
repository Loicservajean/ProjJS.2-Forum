package controllers

import (
	"html/template"
	"net/http"
	"rompelago/dto"
	"rompelago/services"
)

type WebAuthControllers struct {
	authService *services.AuthService
	templates   *template.Template
}

func InitWebAuthController(authService *services.AuthService) *WebAuthControllers {
	fonctionsDisponiblesDansLesTemplates := template.FuncMap{
		"additionner": func(a, b int) int { return a + b },
		"soustraire":  func(a, b int) int { return a - b },
	}
	tmpl := template.Must(template.New("").Funcs(fonctionsDisponiblesDansLesTemplates).ParseGlob("templates/*.html"))
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

	loginResp, err := c.authService.Login(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		c.templates.ExecuteTemplate(w, "connection", map[string]string{
			"Error": err.Error(),
		})
		return
	}

	// Le JWT était généré mais jamais transmis au navigateur : sans ce cookie,
	// le WebAuthMiddleware ne trouve rien sur les pages suivantes et l'utilisateur
	// apparaît déconnecté juste après s'être connecté.
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    loginResp.AccessToken,
		Path:     "/",
		MaxAge:   9 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure devrait être à true en production (HTTPS) ; laissé à false ici
		// pour fonctionner en dev local sur http://localhost.
		Secure: false,
	})

	http.Redirect(w, r, "/forum", http.StatusSeeOther)
}
