package database

import (
	"context"
	"log"
	"os"
)

func DatabaseInit() {
	// database connection
	ctx := context.Background()
	dbUrl := os.Getenv("DATABASE_URL")
	pg, err := NewPG(ctx, dbUrl)
	if err != nil {
		log.Fatalf("Error Creating Database Connection")
	}

	// Ping the Connection
	err = pg.Ping(ctx)
	if err != nil {
		log.Fatalf("Cannot Ping the database: %v", err.Error())
	} else {
		log.Println("Connection Successful")
	}

	defer pg.Close() //close the database

}
