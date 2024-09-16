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
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

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

	// server
	mux := mux.NewRouter()
	sessionStore := auth.NewSessionStore(auth.SessionsOptions{
		CookiesKey: os.Getenv("COOKIES_SECRET_KEY"),
		MaxAge:     86400 * 30,
		Secure:     false,
		HttpOnly:   false,
	})
	authService := auth.NewAuthService(sessionStore) // create a auth service

	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})

	corsHandler := corsOptions.Handler(mux) //cors enabler

	handler := handlers.New(pg, authService)
	mux.Handle("/", http.FileServer(http.Dir("templates/"))) // fix :- need to change to next.js asap

	mux.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Main")
	})
	// auth routes
	mux.HandleFunc("/auth/{provider}/login", handler.HandleProviderLogin).Methods("GET")
	mux.HandleFunc("/auth/{provider}/callback", handler.HandleCallbackFunction).Methods("GET")
	mux.HandleFunc("/auth/logout/{provider}", handler.HandleLogout).Methods("GET")

	log.Fatalln(http.ListenAndServe(PORT, corsHandler)) // run the server
}
