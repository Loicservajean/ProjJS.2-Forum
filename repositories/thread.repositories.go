package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"rompelago/models"
)

type StatusRepositories struct {
	dbContext *sql.DB
}

type CategoryRepositories struct {
	dbContext *sql.DB
}

type FilsRepositories struct {
	dbContext *sql.DB
}

func InitStatusRepositories(dbContext *sql.DB) *StatusRepositories {
	return &StatusRepositories{dbContext: dbContext}
}

func InitCategoryRepositories(dbContext *sql.DB) *CategoryRepositories {
	return &CategoryRepositories{dbContext: dbContext}
}

func InitFilsRepositories(dbContext *sql.DB) *FilsRepositories {
	return &FilsRepositories{dbContext: dbContext}
}

func (r *StatusRepositories) ReadAll() ([]models.StatusModel, error) {
	query := "SELECT id, nom, description FROM statuts;"
	result, resultErr := r.dbContext.Query(query)
	if resultErr != nil {
		return nil, fmt.Errorf("Erreur lors de la requête - %v", resultErr)
	}
	defer result.Close()

	var listStatus []models.StatusModel
	for result.Next() {
		var status models.StatusModel
		scanErr := result.Scan(&status.Id, &status.Nom, &status.Description)
		if scanErr != nil {
			log.Printf("Erreur lors du scan - %v", scanErr)
			continue
		}
		listStatus = append(listStatus, status)
	}
	return listStatus, nil
}

func (r *StatusRepositories) ReadById(id int) (models.StatusModel, error) {
	var status models.StatusModel

	query := "SELECT id, nom, description FROM statuts WHERE id = ?;"
	sqlErr := r.dbContext.QueryRow(query, id).
		Scan(&status.Id, &status.Nom, &status.Description)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return models.StatusModel{}, nil
		}
		return models.StatusModel{}, fmt.Errorf("Erreur lors de la requête - %v", sqlErr)
	}
	return status, nil
}

func (r *CategoryRepositories) ReadAll() ([]models.CategoriesDiscussion, error) {
	query := "SELECT id, nom, description FROM categories;"
	result, resultErr := r.dbContext.Query(query)
	if resultErr != nil {
		return nil, fmt.Errorf("Erreur lors de la requête - %v", resultErr)
	}
	defer result.Close()

	var listCategory []models.CategoriesDiscussion
	for result.Next() {
		var category models.CategoriesDiscussion
		scanErr := result.Scan(&category.Id, &category.Name, &category.Description)
		if scanErr != nil {
			log.Printf("Erreur lors du scan - %v", scanErr)
			continue
		}
		listCategory = append(listCategory, category)
	}
	return listCategory, nil
}

func (r *CategoryRepositories) ReadById(id int) (models.CategoriesDiscussion, error) {
	var category models.CategoriesDiscussion

	query := "SELECT id, nom, description FROM categories WHERE id = ?;"
	sqlErr := r.dbContext.QueryRow(query, id).
		Scan(&category.Id, &category.Name, &category.Description)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return models.CategoriesDiscussion{}, nil
		}
		return models.CategoriesDiscussion{}, fmt.Errorf("Erreur lors de la requête - %v", sqlErr)
	}
	return category, nil
}

func (r *FilsRepositories) ReadAll() ([]models.FilDiscussionModel, error) {
	query := "SELECT id, titre, description, date_creation, date_echeance, statut_id, categorie_id FROM taches;"
	result, resultErr := r.dbContext.Query(query)
	if resultErr != nil {
		return nil, fmt.Errorf("Erreur lors de la requete - %v", resultErr)
	}

	defer result.Close()

	var listTask []models.FilDiscussionModel
	for result.Next() {
		var fildediscussion models.FilDiscussionModel
		scanErr := result.Scan(&fildediscussion.Id, &fildediscussion.Name, &fildediscussion.Description, &fildediscussion.DateCreation, &fildediscussion.StatusId, &fildediscussion.CategorieId)
		if scanErr != nil {
			log.Printf("Erreur lors du scan - %v", scanErr)
			continue
		}
		listTask = append(listTask, fildediscussion)
	}
	return listTask, nil
}

func (r *FilsRepositories) Create(fildediscussion models.FilDiscussionModel) (int, error) {
	query := "INSERT INTO `taches`(`titre`, `description`, `date_creation`, `date_echeance`, `statut_id`, `categorie_id`) VALUES (?,?,?,?,?,?);"

	sqlResult, sqlErr := r.dbContext.Exec(query,
		fildediscussion.Name,
		fildediscussion.Description,
		fildediscussion.DateCreation,
		fildediscussion.StatusId,
		fildediscussion.CategorieId,
	)
	if sqlErr != nil {
		return -1, fmt.Errorf(" Erreur ajout tache - Erreur : \n\t %s", sqlErr.Error())
	}

	id, idErr := sqlResult.LastInsertId()
	if idErr != nil {
		return -1, fmt.Errorf(" Erreur ajout tache - Erreur recuperation identifiant : \n\t %s", idErr.Error())
	}
	return int(id), nil
}

func (r *FilsRepositories) ReadById(id int) (models.FilDiscussionModel, error) {
	var fildediscussion models.FilDiscussionModel

	query := "SELECT id, titre, description, date_creation, date_echeance, statut_id, categorie_id FROM `taches` WHERE `taches`.id = ?;"
	sqlErr := r.dbContext.QueryRow(query, id).
		Scan(&fildediscussion.Id, &fildediscussion.Name, &fildediscussion.Description, &fildediscussion.DateCreation, &fildediscussion.StatusId, &fildediscussion.CategorieId)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return models.FilDiscussionModel{}, nil
		}
		return models.FilDiscussionModel{}, fmt.Errorf(" Erreur lors de la requete - %v", sqlErr)
	}
	return fildediscussion, nil
}

func (r *FilsRepositories) Update(fildediscussion models.FilDiscussionModel) error {
	query := "UPDATE `taches` SET `titre`=?, `description`=?, `date_creation`=?, `date_echeance`=?, `statut_id`=?, `categorie_id`=? WHERE `taches`.id=?;"

	sqlResult, sqlErr := r.dbContext.Exec(query,
		fildediscussion.Name,
		fildediscussion.Description,
		fildediscussion.DateCreation,
		fildediscussion.StatusId,
		fildediscussion.CategorieId,
		fildediscussion.Id,
	)

	if sqlErr != nil {
		return fmt.Errorf("Erreur modification produit - Erreur : \n\t %s", sqlErr.Error())
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur récupération lignes modifiées : %w", err)
	}

	if rowsAffected <= 0 {
		return fmt.Errorf("Erreur modification produit - Aucune ligne modifiée")
	}

	return nil
}

func (r *FilsRepositories) Delete(id int) error {
	sqlResult, sqlErr := r.dbContext.Exec("DELETE FROM `taches` WHERE `taches`.id=?;", id)
	if sqlErr != nil {
		return fmt.Errorf("Erreur suppression produit - Erreur : \n\t %s", sqlErr.Error())
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return fmt.Errorf("Erreur récupération lignes supprimées : %w", err)
	}

	if rowsAffected <= 0 {
		return fmt.Errorf("Erreur suppression produit - Aucune ligne supprimée")
	}

	return nil
}

func (r *FilsRepositories) ReadAllWithCategoryAndStatus() ([]models.FilDiscussionFull, error) {
	query := `
        SELECT 
            t.id, t.titre, t.description, t.date_creation, t.date_echeance,
            t.statut_id, t.categorie_id,
            c.nom, c.description,
            s.nom, s.description
        FROM taches t
        INNER JOIN categories c ON t.categorie_id = c.id
        INNER JOIN statuts s ON t.statut_id = s.id;
    `

	result, err := r.dbContext.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de la requete - %v", err)
	}
	defer result.Close()

	var list []models.FilDiscussionFull

	for result.Next() {
		var t models.FilDiscussionFull
		scanErr := result.Scan(
			&t.Id, &t.Name, &t.Description, &t.DateCreation,
			&t.StatusId, &t.CategorieId,
			&t.CategorieName, &t.CategorieDescription,
			&t.TagName, &t.TagDescription,
		)
		if scanErr != nil {
			log.Printf("Erreur scan - %v", scanErr)
			continue
		}
		list = append(list, t)
	}

	return list, nil
}

func (r *FilsRepositories) ReadByIdWithCategoryAndStatus(id int) (models.FilDiscussionFull, error) {
	query := `
		SELECT 
			t.id, t.titre, t.description, t.date_creation, t.date_echeance,
			t.statut_id, t.categorie_id,
			c.nom, c.description,
			s.nom, s.description
		FROM taches t
		INNER JOIN categories c ON t.categorie_id = c.id
		INNER JOIN statuts s ON t.statut_id = s.id
		WHERE t.id = ?;
	`

	var t models.FilDiscussionFull
	err := r.dbContext.QueryRow(query, id).Scan(
		&t.Id, &t.Name, &t.Description, &t.DateCreation,
		&t.StatusId, &t.CategorieId,
		&t.CategorieName, &t.CategorieDescription,
		&t.TagName, &t.TagDescription,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.FilDiscussionFull{}, fmt.Errorf("tâche introuvable")
		}
		return models.FilDiscussionFull{}, fmt.Errorf("erreur lors de la requête - %v", err)
	}
	return t, nil
}
