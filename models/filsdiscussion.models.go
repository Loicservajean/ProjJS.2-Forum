package models

type FilDiscussionModel struct {
	Id           int
	Name         string
	Description  string
	DateCreation string
	Open         int
	CategorieId  int
	TagId        int
	Creator      Utilisateur
	Archive      int
	CreatorID    int
}

type FilDiscussionWithCategorie struct {
	Id                   int
	Name                 string
	Description          string
	DateCreation         string
	Open                 int
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
	Open                 int
	CategorieId          int
	TagId                int
	Creator              Utilisateur
	CreatorID            int
	Archive              int
	CategorieName        string
	CategorieDescription string
	TagName              string
	TagDescription       string
}
