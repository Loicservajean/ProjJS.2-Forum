package models

type FilDiscussionModel struct {
	Id           int
	Name         string
	Description  string
	DateCreation string
	StatusId     int
	CategorieId  int
	TagId        int
	Creator      Utilisateur
	Archive      int
}

type FilDiscussionWithCategorie struct {
	Id                   int
	Name                 string
	Description          string
	DateCreation         string
	StatusId             int
	CategorieId          int
	TagId                int
	Creator              Utilisateur
	Archive              int
	CategorieName        string
	CategorieDescription string
}

type FilDiscussionFull struct {
	Id                   int
	Name                 string
	Description          string
	DateCreation         string
	StatusId             int
	CategorieId          int
	TagId                int
	Creator              Utilisateur
	Archive              int
	CategorieName        string
	CategorieDescription string
	TagName              string
	TagDescription       string
}
