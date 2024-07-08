package database

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

// Service represents a service that interacts with a database.

type DbService struct {
	Db *gorm.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *DbService
)

func New() *DbService {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", username, password, host, port, database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &DbService{
		Db: db,
	}
	return dbInstance
}

func EnableUuidExtension(db *gorm.DB) error {
	return db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
}
