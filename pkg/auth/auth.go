package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type AuthService struct{}

func NewAuthService(store sessions.Store) *AuthService {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error Loading .Env Files: %v", err)
	}
	gothic.Store = store

	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			os.Getenv("GOOGLE_CALLBACK_URL"),
		),
	)

	return &AuthService{}
}

func (s *AuthService) GetSessionUser(r *http.Request) (goth.User, error) {
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		return goth.User{}, err
	}

	u := session.Values["user"]
	if u == nil {
		return goth.User{}, fmt.Errorf("User is not authenticated! %v", u)
	}
	return u.(goth.User), nil
}

func (s *AuthService) StoreUserSession(w http.ResponseWriter, r *http.Request, user goth.User) error {
	session, _ := gothic.Store.Get(r, SessionName)

	session.Values["user"] = user

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

// remove the user session
func (s *AuthService) RemoveUserSession(w http.ResponseWriter, r *http.Request) {
	// get the request and session
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	session.Values["user"] = goth.User{}

	session.Options.MaxAge = -1

	session.Save(r, w)
}

// middleware
func RequireAuth(handlerFunc http.HandlerFunc, auth *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := auth.GetSessionUser(r)
		if err != nil {
			log.Println("User is not authenticated")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		}

		log.Printf("user is authenticated! user: %v!", session.FirstName)

		handlerFunc(w, r)
	}
}
