package dto

import "rompelago/models"

type FilsResponseDto struct {
    Categories []models.CategoriesDiscussion `json:"categories"`
    Fils       []models.FilDiscussionFull    `json:"fils"`
}
