package main

import (
	"fmt"
	"h-two/internal/server"
)

func main() {
	//dbService := database.New()
	//err := database.EnableUuidExtension(dbService.Db)
	//if err != nil {
	//	log.Fatalf("Failed to enable UUID extension: %v", err)
	//}
	//err = models.Migrate(dbService.Db)
	//if err != nil {
	//	log.Fatalf("Failed to migrate database: %v", err)
	//}
	mainServer := server.NewServer()

	err := mainServer.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
