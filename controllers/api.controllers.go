package controllers

import (
	"net/http"
	"rompelago/dto"
	"rompelago/helper"
	"rompelago/repositories"
	"rompelago/services"
)

type ApiAuthController struct {
	authService *services.AuthService
}

type ApiFilsController struct {
	filService *services.FilDiscussionService
	catRepo    *repositories.CategoryRepositories
}

func InitApiAuthController(authService *services.AuthService) *ApiAuthController {
	return &ApiAuthController{authService: authService}
}

func InitApiFilsController(filService *services.FilDiscussionService, catRepo *repositories.CategoryRepositories) *ApiFilsController {
	return &ApiFilsController{filService: filService, catRepo: catRepo}
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
	categories, err := c.catRepo.ReadAll()
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, "erreur récupération des catégories")
		return
	}
	// (0, 0, 0) = on lit tous les fils sans pagination.
	fils, err := c.filService.ReadAllFull(0, 0, 0)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, "erreur récupération des fils")
		return
	}
	helper.WriteJSON(w, http.StatusOK, dto.FilsResponseDto{
		Categories: categories,
		Fils:       fils,
	})
}
