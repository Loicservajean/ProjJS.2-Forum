package models

type PostModel struct {
	Id         int
	Name       string
	Contenu    string
	DateEnvoi  string
	Creator    Utilisateur
	StatusId   int
	FilAssocié FilDiscussionModel
	ScorePop   int
	NbLike     int
	NbDislike  int
}
