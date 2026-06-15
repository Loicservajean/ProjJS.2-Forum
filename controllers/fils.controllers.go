package controllers

import (
	"encoding/json"
	"net/http"
	"rompelago/helper"
	"rompelago/models"
	"rompelago/services"
	"strconv"

	"github.com/gorilla/mux"
)

type FilsControllers struct {
	service *services.FilDiscussionService
}

func InitFilsController(service *services.FilDiscussionService) *FilsControllers {
	return &FilsControllers{service: service}
}

func readFilId(r *http.Request) (int, error) {
	return strconv.Atoi(mux.Vars(r)["id"])
}

func (c *FilsControllers) ReadAll(w http.ResponseWriter, r *http.Request) {
	fils, err := c.service.ReadAll()
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSON(w, http.StatusOK, fils)
}

func (c *FilsControllers) Create(w http.ResponseWriter, r *http.Request) {
	var fil models.FilDiscussionFull
	if err := json.NewDecoder(r.Body).Decode(&fil); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "JSON invalide")
		return
	}
	id, err := c.service.Create(fil)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	helper.WriteJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func (c *FilsControllers) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := readFilId(r)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Identifiant invalide")
		return
	}
	if err := c.service.Delete(id); err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Tâche supprimée"})
}

func (c *FilsControllers) ReadByIdDetails(w http.ResponseWriter, r *http.Request) {
	id, err := readFilId(r)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Identifiant invalide")
		return
	}
	fil, err := c.service.ReadByIdWithDetails(id)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSON(w, http.StatusOK, fil)
}
