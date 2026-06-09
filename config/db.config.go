package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func InitDbContext() *sql.DB {
	dbUser := GetRequiredEnv("DB_USER")
	DbPwd := GetRequiredEnv("DB_PWD")
	dbHost := GetRequiredEnv("DB_HOST")
	dbPort := GetRequiredEnv("DB_PORT")
	dbName := GetRequiredEnv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, DbPwd, dbHost, dbPort, dbName)

	dbContext, dbContextErr := sql.Open("mysql", dsn)
	if dbContextErr != nil {
		log.Fatalf("Erreur base de donnée - Initialisation de la connection impossible \n %v", dbContextErr)
	}

	pingErr := dbContext.Ping()
	if pingErr != nil {
		log.Fatalf("Erreur base de donnée - Ping impossible \n %v", pingErr)
	}
	log.Printf("log base de donnée - connexion reussie")
	return dbContext
}
