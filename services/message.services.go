package services

import (
	"fmt"
	"rompelago/models"
	"rompelago/repositories"
	"time"
)

type PostDiscussionService struct {
	PostRepository *repositories.PostRepositories
}

func InitPostDiscussionService(repo *repositories.PostRepositories) *PostDiscussionService {
	return &PostDiscussionService{PostRepository: repo}
}

func (s *PostDiscussionService) ReadByFilId(filId int) ([]models.PostModel, error) {
	if filId <= 0 {
		return nil, fmt.Errorf("identifiant de fil invalide : %d", filId)
	}
	return s.PostRepository.ReadPostsByFilId(filId)
}

func (s *PostDiscussionService) Create(post models.PostModel) (int, error) {
	if post.Contenu == "" {
		return -1, fmt.Errorf("le contenu est obligatoire")
	}
	if post.Creator.Id <= 0 {
		return -1, fmt.Errorf("utilisateur invalide")
	}
	if post.FilAssocié.Id <= 0 {
		return -1, fmt.Errorf("fil de discussion invalide")
	}
	post.DateEnvoi = time.Now().Format("2006-01-02 15:04:05")
	return s.PostRepository.CreatePost(post)
}

func (s *PostDiscussionService) Delete(id int) error {
	if id <= 0 {
		return fmt.Errorf("identifiant invalide : %d", id)
	}
	return s.PostRepository.DeletePost(id)
}

func (s *PostDiscussionService) LikeDislike(userId int, messageId int, action string) error {
	if messageId <= 0 || userId <= 0 {
		return fmt.Errorf("identifiant invalide")
	}
	if action != "like" && action != "dislike" {
		return fmt.Errorf("action invalide")
	}

	// Vérifier le vote existant
	existingVote, err := s.PostRepository.GetLikeDislike(userId, messageId)
	if err != nil {
		return err
	}

	// Déjà voté pareil → on ignore
	if existingVote == action {
		return fmt.Errorf("vous avez déjà %s ce message", action)
	}

	// Annuler l'ancien vote opposé si besoin
	if existingVote != "" {
		if existingVote == "like" {
			s.PostRepository.UpdateLikeDislike(messageId, "cancel_like")
		} else {
			s.PostRepository.UpdateLikeDislike(messageId, "cancel_dislike")
		}
	}

	// Appliquer le nouveau vote
	if err := s.PostRepository.UpdateLikeDislike(messageId, action); err != nil {
		return err
	}

	// Sauvegarder en BDD
	return s.PostRepository.UpsertLikeDislike(userId, messageId, action)
}
