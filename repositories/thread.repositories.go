package repositories

import "database/sql"

type StatusRepositories struct {
	dbContext *sql.DB
}

type CategoryRepositories struct {
	dbContext *sql.DB
}
