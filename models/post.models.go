package models

type PostModel struct {
	Id           int
	Name         string
	Contenu      string
	DateCreation string
	Creator      Utilisateur
	StatusId     int
	FilAssocié   FilDiscussionModel
	ScorePop     int
	NbLike       int
	NbDislike    int
}

type PostWithCategorie struct {
	Id                   int
	Name                 string
	Contenu              string
	DateCreation         string
	Creator              Utilisateur
	StatusId             int
	FilAssocié           FilDiscussionModel
	ScorePop             int
	NbLike               int
	NbDislike            int
	CategorieName        string
	CategorieDescription string
}

type Post struct {
	Id                   int
	Name                 string
	Contenu              string
	DateCreation         string
	Creator              Utilisateur
	StatusId             int
	FilAssocié           FilDiscussionModel
	ScorePop             int
	NbLike               int
	NbDislike            int
	CategorieDescription string
	TagName              string
	TagDescription       string
}
