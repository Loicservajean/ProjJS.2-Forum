package models

type Utilisateur struct {
	Id        int
	Name      string
	email     string
	Favoris   []FilDiscussionModel
	StatusBan int
}
