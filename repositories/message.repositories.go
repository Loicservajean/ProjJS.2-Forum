package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"rompelago/models"
)

type PostRepositories struct {
	dbContext *sql.DB
}

func InitPostRepositories(dbContext *sql.DB) *PostRepositories {
	return &PostRepositories{dbContext: dbContext}
}

// limit à 0 veut dire qu'on ne filtre rien du tout.
func (r *PostRepositories) ReadPostsByFilId(filId int, limit, offset int, sort string) ([]models.PostModel, error) {
	orderBy := "m.date_envoi DESC"
	if sort == "popularity" {
		orderBy = "m.scorepop DESC, m.date_envoi DESC"
	}
	query := `
		SELECT m.id_message, m.name, m.contenu, m.date_envoi, m.scorepop, m.nb_like, m.nb_dislike,
		       u.id_utilisateur, u.pseudo
		FROM Message m
		LEFT JOIN Utilisateur u ON u.id_utilisateur = m.fk_utilisateur
		WHERE m.fk_fil_de_discussion = ?
		ORDER BY ` + orderBy + `
	`
	args := []interface{}{filId}
	if limit > 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, limit, offset)
	}

	result, err := r.dbContext.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture messages du fil %d - %v", filId, err)
	}
	defer result.Close()

	var list []models.PostModel
	for result.Next() {
		var post models.PostModel
		scanErr := result.Scan(
			&post.Id, &post.Name, &post.Contenu, &post.DateEnvoi,
			&post.ScorePop, &post.NbLike, &post.NbDislike,
			&post.Creator.Id, &post.Creator.Name,
		)
		if scanErr != nil {
			log.Printf("erreur scan message - %v", scanErr)
			continue
		}
		list = append(list, post)
	}
	return list, nil
}

// CountMessagesByFilId va juste compter combien de messages traînent dans un fil donné.
func (r *PostRepositories) CountMessagesByFilId(filId int) (int, error) {
	var total int
	err := r.dbContext.QueryRow("SELECT COUNT(*) FROM Message WHERE fk_fil_de_discussion = ?;", filId).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("erreur comptage messages du fil %d - %v", filId, err)
	}
	return total, nil
}

func (r *PostRepositories) ReadAllMessages() ([]models.PostModel, error) {
	query := `
		SELECT m.id_message, m.name, m.contenu, m.date_envoi, m.scorepop, m.nb_like, m.nb_dislike,
		       m.fk_fil_de_discussion, u.id_utilisateur, u.pseudo
		FROM Message m
		LEFT JOIN Utilisateur u ON u.id_utilisateur = m.fk_utilisateur
		ORDER BY m.date_envoi DESC
	`

	result, err := r.dbContext.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture messages admin - %v", err)
	}
	defer result.Close()

	var list []models.PostModel
	for result.Next() {
		var post models.PostModel
		scanErr := result.Scan(
			&post.Id, &post.Name, &post.Contenu, &post.DateEnvoi,
			&post.ScorePop, &post.NbLike, &post.NbDislike,
			&post.FilAssocié.Id, &post.Creator.Id, &post.Creator.Name,
		)
		if scanErr != nil {
			log.Printf("erreur scan message admin - %v", scanErr)
			continue
		}
		list = append(list, post)
	}
	return list, nil
}

// CreatePost insère un message dans un fil de discussion.
func (r *PostRepositories) CreatePost(post models.PostModel) (int, error) {
	query := `
		INSERT INTO Message (fk_fil_de_discussion, fk_utilisateur, contenu, name, date_envoi, nb_like, nb_dislike, scorepop)
		VALUES (?, ?, ?, ?, ?, 0, 0, 0);
	`
	sqlResult, sqlErr := r.dbContext.Exec(query,
		post.FilAssocié.Id,
		post.Creator.Id,
		post.Contenu,
		post.Name,
		post.DateEnvoi,
	)
	if sqlErr != nil {
		return -1, fmt.Errorf("erreur ajout message - %s", sqlErr.Error())
	}

	id, idErr := sqlResult.LastInsertId()
	if idErr != nil {
		return -1, fmt.Errorf("erreur récupération identifiant message - %s", idErr.Error())
	}
	return int(id), nil
}

// ReadPostById retourne un message par son id.
func (r *PostRepositories) ReadPostById(id int) (models.PostModel, error) {
	var post models.PostModel
	query := `
		SELECT m.id_message, m.contenu, m.date_envoi, m.scorepop, m.nb_like, m.nb_dislike,
		       u.id_utilisateur, u.pseudo
		FROM Message m
		LEFT JOIN Utilisateur u ON u.id_utilisateur = m.fk_utilisateur
		WHERE m.id_message = ?;
	`
	err := r.dbContext.QueryRow(query, id).Scan(
		&post.Id, &post.Contenu, &post.DateEnvoi,
		&post.ScorePop, &post.NbLike, &post.NbDislike,
		&post.Creator.Id, &post.Creator.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.PostModel{}, nil
		}
		return models.PostModel{}, fmt.Errorf("erreur lecture message %d - %v", id, err)
	}
	return post, nil
}

// DeletePost supprime un message par son id (et ses votes associés).
func (r *PostRepositories) DeletePost(id int) error {
	_, sqlErr := r.dbContext.Exec("DELETE FROM LikeDislike WHERE fk_message = ?;", id)
	if sqlErr != nil {
		return fmt.Errorf("erreur suppression votes du message - %s", sqlErr.Error())
	}

	sqlResult, sqlErr := r.dbContext.Exec("DELETE FROM Message WHERE id_message = ?;", id)
	if sqlErr != nil {
		return fmt.Errorf("erreur suppression message - %s", sqlErr.Error())
	}
	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur récupération lignes supprimées : %w", err)
	}
	if rowsAffected <= 0 {
		return fmt.Errorf("aucun message supprimé (id=%d)", id)
	}
	return nil
}

// Récupère le vote existant d'un utilisateur sur un message
func (r *PostRepositories) GetLikeDislike(userId int, messageId int) (string, error) {
	var typeVote string
	err := r.dbContext.QueryRow(
		`SELECT type_vote FROM LikeDislike WHERE fk_utilisateur = ? AND fk_message = ?`,
		userId, messageId,
	).Scan(&typeVote)

	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("erreur lecture vote - %v", err)
	}
	return typeVote, nil
}

// Insère ou met à jour le vote
func (r *PostRepositories) UpsertLikeDislike(userId int, messageId int, action string) error {
	_, err := r.dbContext.Exec(
		`INSERT INTO LikeDislike (fk_utilisateur, fk_message, type_vote) VALUES (?, ?, ?)
         ON DUPLICATE KEY UPDATE type_vote = ?`,
		userId, messageId, action, action,
	)
	return err
}

// UpdatePost met à jour le titre et le contenu d'un message.
func (r *PostRepositories) UpdatePost(id int, name, contenu string) error {
	query := `UPDATE Message SET name = ?, contenu = ? WHERE id_message = ?;`
	_, sqlErr := r.dbContext.Exec(query, name, contenu, id)
	if sqlErr != nil {
		return fmt.Errorf("erreur modification message - %s", sqlErr.Error())
	}
	return nil
}

// ReadPostByIdFull retourne un message par son id avec le nom du fil associé.
func (r *PostRepositories) ReadPostByIdFull(id int) (models.PostModel, error) {
	var post models.PostModel
	query := `
		SELECT m.id_message, m.name, m.contenu, m.date_envoi, m.scorepop, m.nb_like, m.nb_dislike,
		       u.id_utilisateur, u.pseudo, m.fk_fil_de_discussion
		FROM Message m
		LEFT JOIN Utilisateur u ON u.id_utilisateur = m.fk_utilisateur
		WHERE m.id_message = ?;
	`
	err := r.dbContext.QueryRow(query, id).Scan(
		&post.Id, &post.Name, &post.Contenu, &post.DateEnvoi,
		&post.ScorePop, &post.NbLike, &post.NbDislike,
		&post.Creator.Id, &post.Creator.Name,
		&post.FilAssocié.Id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.PostModel{}, fmt.Errorf("message introuvable")
		}
		return models.PostModel{}, fmt.Errorf("erreur lecture message %d - %v", id, err)
	}
	return post, nil
}
func (r *PostRepositories) UpdateLikeDislike(messageId int, action string) error {
	var query string
	switch action {
	case "like":
		query = `UPDATE Message SET nb_like = nb_like + 1, scorepop = scorepop + 1 WHERE id_message = ?`
	case "dislike":
		query = `UPDATE Message SET nb_dislike = nb_dislike + 1, scorepop = scorepop - 1 WHERE id_message = ?`
	case "cancel_like":
		query = `UPDATE Message SET nb_like = nb_like - 1, scorepop = scorepop - 1 WHERE id_message = ?`
	case "cancel_dislike":
		query = `UPDATE Message SET nb_dislike = nb_dislike - 1, scorepop = scorepop + 1 WHERE id_message = ?`
	default:
		return fmt.Errorf("action invalide : %s", action)
	}

	result, err := r.dbContext.Exec(query, messageId)
	if err != nil {
		return fmt.Errorf("erreur mise à jour vote - %v", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("aucun message trouvé avec l'id %d", messageId)
	}
	return nil
}
