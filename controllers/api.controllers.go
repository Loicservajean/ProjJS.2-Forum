package controllers

import (
	"net/http"
	"rompelago/dto"
	"rompelago/helper"
	"rompelago/services"
)

type ApiAuthController struct {
	authService *services.AuthService
}

type ApiFilsController struct {
	filService *services.FilDiscussionService
}

func InitApiAuthController(authService *services.AuthService) *ApiAuthController {
	return &ApiAuthController{authService: authService}
}

func InitApiFilsController(filService *services.FilDiscussionService) *ApiFilsController {
	return &ApiFilsController{filService: filService}
}

// POST /api/auth/login
func (c *ApiAuthController) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	req := dto.LoginRequestDto{
		Pseudo:   r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if r.FormValue("confirm_password") != req.Password {
		helper.WriteError(w, http.StatusBadRequest, "les mots de passe ne correspondent pas")
		return
	}

	resp, err := c.authService.Login(req)
	if err != nil {
		helper.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, resp)
}

// POST /api/auth/register
func (c *ApiAuthController) Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	req := dto.RegisterRequestDto{
		Pseudo:   r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if r.FormValue("confirm_password") != req.Password {
		helper.WriteError(w, http.StatusBadRequest, "les mots de passe ne correspondent pas")
		return
	}

	if err := c.authService.Register(req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusCreated, dto.ResponseDto{
		Code:    http.StatusCreated,
		Message: "compte créé avec succès",
	})
}

// GET /api/fils
func (c *ApiFilsController) ListFils(w http.ResponseWriter, r *http.Request) {
	// (0, 0, 0) = on lit tous les fils sans pagination.
	fils, err := c.filService.ReadAllFull(0, 0, 0)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, "erreur récupération des fils")
		return
	}
	helper.WriteJSON(w, http.StatusOK, fils)
}
