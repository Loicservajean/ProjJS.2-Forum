package services

import (
	"fmt"
	"rompelago/models"
	"rompelago/repositories"
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

func (s *FilDiscussionService) ReadAllFull() ([]models.FilDiscussionFull, error) {
	return s.FilsRepository.ReadAllWithCategoryAndStatus()
}

func (s *FilDiscussionService) Create(fils models.FilDiscussionFull) (int, error) {
	if fils.Name == "" || fils.Description == "" {
		return -1, fmt.Errorf("titre et description obligatoires")
	}
	return s.FilsRepository.Create(models.FilDiscussionModel{
		Id:          fils.Id,
		Name:        fils.Name,
		Description: fils.Description,
		Open:        fils.Open,
		Creator:     fils.Creator,
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
