package services

import (
	"fmt"
	"rompelago/models"
	"rompelago/repositories"
	"time"
)

type FilDiscussionService struct {
	FilsRepository *repositories.FilsRepositories
}

func InitFilDiscussionService(repo *repositories.FilsRepositories) *FilDiscussionService {
	return &FilDiscussionService{FilsRepository: repo}
}

func (s *FilDiscussionService) ReadAll() ([]models.FilDiscussionModel, error) {
	return s.FilsRepository.ReadAll()
}

func (s *FilDiscussionService) ReadAllFull(limit, offset int) ([]models.FilDiscussionFull, error) {
	return s.FilsRepository.ReadAllWithCategoryAndStatus(limit, offset)
}

func (s *FilDiscussionService) CountFils() (int, error) {
	return s.FilsRepository.CountFils()
}

func (s *FilDiscussionService) Create(fils models.FilDiscussionFull) (int, error) {
	if fils.Name == "" || fils.Description == "" {
		return -1, fmt.Errorf("titre et description obligatoires")
	}
	dateNow := time.Now().Format("2006-01-02 15:04:05")
	return s.FilsRepository.Create(models.FilDiscussionModel{
		Id:           fils.Id,
		Name:         fils.Name,
		Description:  fils.Description,
		DateCreation: dateNow,
		Open:         fils.Open,
		CreatorID:    fils.CreatorID,
	})
}

func (s *FilDiscussionService) Delete(id int) error {
	if id <= 0 {
		return fmt.Errorf("identifiant invalide : %d", id)
	}
	return s.FilsRepository.Delete(id)
}

func (s *FilDiscussionService) ReadByIdWithDetails(id int) (models.FilDiscussionFull, error) {
	if id <= 0 {
		return models.FilDiscussionFull{}, fmt.Errorf("identifiant invalide : %d", id)
	}
	return s.FilsRepository.ReadByIdWithCategoryAndStatus(id)
}

func (s *FilDiscussionService) UpdateFull(id int, fils models.FilDiscussionFull) error {
	if id <= 0 {
		return fmt.Errorf("identifiant invalide : %d", id)
	}
	if fils.Description == "" {
		return fmt.Errorf("description obligatoire")
	}
	return s.FilsRepository.UpdateFull(models.FilDiscussionModel{
		Id:          id,
		Name:        fils.Name,
		Description: fils.Description,
		DateCreation: fils.DateCreation,
		Open:        fils.Open,
		Archive:     fils.Archive,
		CreatorID:   fils.CreatorID,
	})
}
