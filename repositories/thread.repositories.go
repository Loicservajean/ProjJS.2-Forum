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
	query := "SELECT id_status, name, Description FROM Status;"
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

	query := "SELECT id_status, name, Description FROM Status WHERE id_status = ?;"
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
	query := "SELECT id_type, name, Description FROM CategoriesDiscussion;"
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

	query := "SELECT id_type, name, Description FROM CategoriesDiscussion WHERE id_type = ?;"
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
	query := "SELECT id_fil_de_discussion, name, description, date_creation, open, archive FROM Fil_de_discussion;"
	result, resultErr := r.dbContext.Query(query)
	if resultErr != nil {
		return nil, fmt.Errorf("Erreur lors de la requete - %v", resultErr)
	}

	defer result.Close()

	var listFils []models.FilDiscussionModel
	for result.Next() {
		var fildediscussion models.FilDiscussionModel
		scanErr := result.Scan(&fildediscussion.Id, &fildediscussion.Name, &fildediscussion.Description, &fildediscussion.DateCreation, &fildediscussion.Open, &fildediscussion.Archive)
		if scanErr != nil {
			log.Printf("Erreur lors du scan - %v", scanErr)
			continue
		}
		listFils = append(listFils, fildediscussion)
	}
	return listFils, nil
}

func (r *FilsRepositories) Create(fildediscussion models.FilDiscussionModel) (int, error) {
	query := "INSERT INTO `Fil_de_discussion`(`name`, `description`, `date_creation`, `open`, `archive`, `fk_utilisateur`) VALUES (?,?,?,?,?,?);"

	sqlResult, sqlErr := r.dbContext.Exec(query,
		fildediscussion.Name,
		fildediscussion.Description,
		fildediscussion.DateCreation,
		fildediscussion.Open,
		fildediscussion.Archive,
		fildediscussion.CreatorID,
	)
	if sqlErr != nil {
		return -1, fmt.Errorf(" Erreur ajout Fil - Erreur : \n\t %s", sqlErr.Error())
	}

	id, idErr := sqlResult.LastInsertId()
	if idErr != nil {
		return -1, fmt.Errorf(" Erreur ajout Fil - Erreur recuperation identifiant : \n\t %s", idErr.Error())
	}
	return int(id), nil
}

func (r *FilsRepositories) ReadById(id int) (models.FilDiscussionModel, error) {
	var fildediscussion models.FilDiscussionModel

	query := "SELECT id_fil_de_discussion, name, description, date_creation, open, archive FROM `Fil_de_discussion` WHERE `Fil_de_discussion`.id_fil_de_discussion = ?;"
	sqlErr := r.dbContext.QueryRow(query, id).
		Scan(&fildediscussion.Id, &fildediscussion.Name, &fildediscussion.Description, &fildediscussion.DateCreation, &fildediscussion.Open, &fildediscussion.Archive)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return models.FilDiscussionModel{}, nil
		}
		return models.FilDiscussionModel{}, fmt.Errorf(" Erreur lors de la requete - %v", sqlErr)
	}
	return fildediscussion, nil
}

func (r *FilsRepositories) Update(fildediscussion models.FilDiscussionModel) error {
	query := "UPDATE `Fil_de_discussion` SET `name`=?, `description`=?, `date_creation`=?, `open`=?, `archive`=? WHERE `Fil_de_discussion`.id_fil_de_discussion=?;"

	sqlResult, sqlErr := r.dbContext.Exec(query,
		fildediscussion.Name,
		fildediscussion.Description,
		fildediscussion.DateCreation,
		fildediscussion.Open,
		fildediscussion.Archive,
		fildediscussion.Id,
	)

	if sqlErr != nil {
		return fmt.Errorf("Erreur modification fil - Erreur : \n\t %s", sqlErr.Error())
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur récupération lignes modifiées : %w", err)
	}

	if rowsAffected <= 0 {
		return fmt.Errorf("Erreur modification fil - Aucune ligne modifiée")
	}

	return nil
}

func (r *FilsRepositories) Delete(id int) error {
	sqlResult, sqlErr := r.dbContext.Exec("DELETE FROM `Fil_de_discussion` WHERE `Fil_de_discussion`.id_fil_de_discussion=?;", id)
	if sqlErr != nil {
		return fmt.Errorf("Erreur suppression fil - Erreur : \n\t %s", sqlErr.Error())
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return fmt.Errorf("Erreur récupération lignes supprimées : %w", err)
	}

	if rowsAffected <= 0 {
		return fmt.Errorf("Erreur suppression fil - Aucune ligne supprimée")
	}

	return nil
}

func (r *FilsRepositories) ReadAllWithCategoryAndStatus() ([]models.FilDiscussionFull, error) {
	query := `
        SELECT 
            t.id_fil_de_discussion, t.name, t.description, t.date_creation,
            t.open, t.archive,
            COALESCE(c.name, ''), COALESCE(c.Description, ''),
            COALESCE(s.name, ''), COALESCE(s.Description, ''),
            COALESCE(t.fk_utilisateur, 0)
        FROM Fil_de_discussion t
        LEFT JOIN Fil_Type ft ON ft.fk_fil = t.id_fil_de_discussion
		LEFT JOIN CategoriesDiscussion c ON c.id_type = ft.fk_type
		LEFT JOIN Fil_status fs ON fs.fk_fil = t.id_fil_de_discussion
		LEFT JOIN Status s ON s.id_status = fs.fk_status
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
			&t.Open, &t.Archive,
			&t.CategorieName, &t.CategorieDescription,
			&t.TagName, &t.TagDescription,
			&t.Creator.Id,
		)
		if scanErr != nil {
			log.Printf("Erreur scan - %v", scanErr)
			continue
		}
		list = append(list, t)
	}

	return list, nil
}

//Doc pour coalesce en dessous : Ceci est une fonction SQL pour Ignorer un Paramètre qui est NULL dans la Base de données (D'après ce que j'ai compris de la doc)
/*https://learnsql.fr/blog/comment-utiliser-la-fonction-coalesce-en-sql/*/
func (r *FilsRepositories) ReadByIdWithCategoryAndStatus(id int) (models.FilDiscussionFull, error) {
	query := `
		SELECT 
			t.id_fil_de_discussion, t.name, t.description, t.date_creation,
			t.open, t.archive,
			COALESCE(c.name, ''), COALESCE(c.Description, ''),
			COALESCE(s.name, ''), COALESCE(s.Description, ''),
			COALESCE(t.fk_utilisateur, 0)
		FROM Fil_de_discussion t
		LEFT JOIN Fil_Type ft ON ft.fk_fil = t.id_fil_de_discussion
		LEFT JOIN CategoriesDiscussion c ON c.id_type = ft.fk_type
		LEFT JOIN Fil_status fs ON fs.fk_fil = t.id_fil_de_discussion
		LEFT JOIN Status s ON s.id_status = fs.fk_status
		WHERE t.id_fil_de_discussion = ?;
	`

	var t models.FilDiscussionFull
	err := r.dbContext.QueryRow(query, id).Scan(
		&t.Id, &t.Name, &t.Description, &t.DateCreation,
		&t.Open, &t.Archive,
		&t.CategorieName, &t.CategorieDescription,
		&t.TagName, &t.TagDescription,
		&t.Creator.Id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.FilDiscussionFull{}, fmt.Errorf("fil introuvable")
		}
		return models.FilDiscussionFull{}, fmt.Errorf("erreur lors de la requête - %v", err)
	}
	return t, nil
}

func (r *FilsRepositories) UpdateFull(fildediscussion models.FilDiscussionModel) error {
	query := `
		UPDATE Fil_de_discussion
		SET description = ?, open = ?, archive = ?
		WHERE id_fil_de_discussion = ?;
	`

	sqlResult, sqlErr := r.dbContext.Exec(query,
		fildediscussion.Description,
		fildediscussion.Open,
		fildediscussion.Archive,
		fildediscussion.Id,
	)

	if sqlErr != nil {
		return fmt.Errorf("erreur modification fil - %s", sqlErr.Error())
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur récupération lignes modifiées : %w", err)
	}

	if rowsAffected <= 0 {
		return fmt.Errorf("aucune ligne modifiée")
	}

	return nil
}
