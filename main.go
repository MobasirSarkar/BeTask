package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MobasirSarkar/BeTask/database"
	"github.com/MobasirSarkar/BeTask/pkg/auth"
	"github.com/MobasirSarkar/BeTask/pkg/handlers"
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

	// database connection
	ctx := context.Background()
	dbUrl := os.Getenv("DATABASE_URL")
	pg, err := database.NewPG(ctx, dbUrl)
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

	// authentication
	sessionStore := auth.NewCookieStore(auth.SessionOptions{
		CookiesKey: os.Getenv("COOKIES_SECRET_KEY"),
		MaxAge:     60 * 60 * 24 * 30,
		Secure:     false,
		HttpOnly:   false,
	})

	authService := auth.NewAuthService(sessionStore)

	handler := handlers.New(pg, authService)

	// server
	fmt.Printf("Listening to Localhost %v", PORT)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", hello)

	server := http.Server{
		Addr:    PORT,
		Handler: mux,
	}

	mux.HandleFunc("/", hello)
	//auth
	mux.HandleFunc("GET /auth/:provider", handler.HandleProviderLogin)
	mux.HandleFunc("GET /auth/:provider/callback", handler.HandleAuthCallbackFunction)
	mux.HandleFunc("GET /auth/logout/:provider", handler.HandleLogout)

	server.ListenAndServe()
}
