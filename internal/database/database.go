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
	db *gorm.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	xataApiKey = os.Getenv("XATA_API_KEY")
	dbInstance *DbService
)

func New() *DbService {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	dsn := fmt.Sprintf("postgresql://d4g40h:%s@us-east-1.sql.xata.sh/h-two:main?sslmode=require", xataApiKey)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &DbService{
		db: db,
	}
	return dbInstance
}
