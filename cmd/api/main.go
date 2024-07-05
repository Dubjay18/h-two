package main

import (
	"fmt"
	"h-two/internal/database"
	"h-two/internal/models"
	"h-two/internal/server"
	"log"
)

func main() {
	dbService := database.New()
	err := database.EnableUuidExtension(dbService.Db)
	if err != nil {
		log.Fatalf("Failed to enable UUID extension: %v", err)
	}
	err = models.Migrate(dbService.Db)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	server := server.NewServer()

	err = server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
