package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MobasirSarkar/BeTask/database"
	"github.com/joho/godotenv"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Mobasir")
}

const PORT = ":8080"

func main() {
	//  helps to load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error Reading .env file: %v", err)
	}

	//database

	database.DatabaseInit()

	// server
	fmt.Printf("Listening to Localhost %v", PORT)
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)

	server := http.Server{
		Addr:    PORT,
		Handler: mux,
	}

	server.ListenAndServe()
}
