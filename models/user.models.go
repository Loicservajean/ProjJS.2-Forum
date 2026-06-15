package models

type Utilisateur struct {
	Id          int
	Name        string
	Email       string
	Description string
	Favoris     []FilDiscussionModel
	StatusBan   int
	Role        string
}
