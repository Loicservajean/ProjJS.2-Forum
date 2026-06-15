package repositories

import (
	"database/sql"
	"fmt"
	"rompelago/models"
)

type UserRepositories struct {
	dbContext *sql.DB
}

func InitUserRepositories(dbContext *sql.DB) *UserRepositories {
	return &UserRepositories{dbContext: dbContext}
}

// FindByPseudo recherche un utilisateur par son pseudo avec son mot de passe hashé
func (r *UserRepositories) FindByPseudo(pseudo string) (models.Utilisateur, string, error) {
	query := `
		SELECT u.id_utilisateur, u.pseudo, u.e_mail, u.description, u.ban, u.status,
		       m.mot_de_passe
		FROM Utilisateur u
		INNER JOIN Mot_de_passe m ON m.fk_utilisateur = u.id_utilisateur
		WHERE u.pseudo = ?
	`
	var user models.Utilisateur
	var hashedPassword string

	err := r.dbContext.QueryRow(query, pseudo).Scan(
		&user.Id, &user.Name, &user.Email, &user.Description, &user.StatusBan, &user.Role,
		&hashedPassword,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Utilisateur{}, "", fmt.Errorf("utilisateur introuvable")
		}
		return models.Utilisateur{}, "", fmt.Errorf("erreur lors de la requête - %v", err)
	}
	return user, hashedPassword, nil
}

// Create insère un nouvel utilisateur et son mot de passe hashé
func (r *UserRepositories) Create(pseudo, email, hashedPassword string) (int, error) {
	// Démarrage d'une transaction :
	// Une transaction c'est quand on veut faire plusieurs requêtes et que si y en a une qui n'est pas bonne, on annule tout
	// Exemple ici : on veut créer un utilisateur et son mot de passe
	// Si on arrive a avoir notre utilisateur dans la DB mais que le mot de passe ne passe pas
	// Alors on supprime l'utilisateur pour pas avoir un utilisateur sans mot de passe

	//Voilà pour les explications maison pour ceux qui ne comprennent pas la ligne en dessous pour qu'ils disent pas que c'est de L'ia
	//N'est ce pas Loic ?
	tx, err := r.dbContext.Begin()
	if err != nil {
		return -1, fmt.Errorf("erreur début transaction - %v", err)
	}
	defer tx.Rollback()

	// Insertion de l'utilisateur (fk_fil_de_discussion vaut 0 par défaut tant que l'utilisateur n'a pas de fil propre)
	res, err := tx.Exec(
		`INSERT INTO Utilisateur (pseudo, e_mail, description, ban, status, fk_fil_de_discussion)
		 VALUES (?, ?, '', FALSE, 'user', NULL)`,
		pseudo, email,
	)
	if err != nil {
		return -1, fmt.Errorf("erreur insertion utilisateur - %v", err)
	}

	userID, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("erreur récupération id utilisateur - %v", err)
	}

	// Insertion du mot de passe
	_, err = tx.Exec(
		`INSERT INTO Mot_de_passe (fk_utilisateur, mot_de_passe) VALUES (?, ?)`,
		userID, hashedPassword,
	)
	if err != nil {
		return -1, fmt.Errorf("erreur insertion mot de passe - %v", err)
	}

	if err := tx.Commit(); err != nil {
		return -1, fmt.Errorf("erreur commit transaction - %v", err)
	}
	return int(userID), nil
}

// FindById recherche un utilisateur par son identifiant.
func (r *UserRepositories) FindById(id int) (models.Utilisateur, error) {
	query := `
		SELECT id_utilisateur, pseudo, e_mail, description, ban, status
		FROM Utilisateur
		WHERE id_utilisateur = ?
	`
	var user models.Utilisateur
	err := r.dbContext.QueryRow(query, id).Scan(
		&user.Id, &user.Name, &user.Email, &user.Description, &user.StatusBan, &user.Role,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Utilisateur{}, fmt.Errorf("utilisateur introuvable")
		}
		return models.Utilisateur{}, fmt.Errorf("erreur lors de la requête - %v", err)
	}
	return user, nil
}

// ExistsByPseudoOrEmail vérifie si un pseudo ou un email est déjà utilisé.
func (r *UserRepositories) ExistsByPseudoOrEmail(pseudo, email string) (bool, error) {
	var count int
	err := r.dbContext.QueryRow(
		`SELECT COUNT(*) FROM Utilisateur WHERE pseudo = ? OR e_mail = ?`,
		pseudo, email,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erreur vérification existence - %v", err)
	}
	return count > 0, nil
}
