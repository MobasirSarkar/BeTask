package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

// handles the login
func (h *Handler) HandleProviderLogin(w http.ResponseWriter, r *http.Request) {
	if u, err := gothic.CompleteUserAuth(w, r); err == nil {
		log.Printf("User already authenticated! %v", u)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		gothic.BeginAuthHandler(w, r)
	}

}

func (h *Handler) HandleCallbackFunction(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("Authenticated User: %v", user.FirstName) // have to fix , name is not visible

	exists, err := h.userExists(user.Email)
	if err != nil {
		log.Printf("Error checking user existence: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} // checks if the is already in the database or not

	err = h.auth.StoreUserSession(w, r, user)
	if err != nil {
		log.Println(err)
		return
	} // store the users sessions

	if !exists {
		err = h.createUser(user)
		if err != nil {
			log.Printf("Error Creating User: %v\n", err)
			http.Error(w, "Internal Server error", http.StatusInternalServerError)
			return
		}
	}
	log.Println(user.Name)
	w.Header().Set("Location", "/main")
	w.WriteHeader(http.StatusTemporaryRedirect)

}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	log.Println("Logging Out....")

	err := gothic.Logout(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	h.auth.RemoveUserSession(w, r)

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) createUser(gothUser goth.User) error {
	uuid := generateUUID()
	query := `INSERT INTO users (id,google_id,profile_pic_url,name,email,created_at,updated_at)
             VALUES ($1,$2,$3,$4,$5,$6,$7)
   `
	_, err := h.db.Getpool().Exec(context.Background(), query, uuid, gothUser.UserID, gothUser.AvatarURL, gothUser.Name, gothUser.Email, time.Now(), time.Now())
	return err
}

func generateUUID() string { // generate random uuid
	return fmt.Sprintf("%s", uuid.NewString())
}

func (h *Handler) userExists(email string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := h.db.Getpool().QueryRow(context.Background(), query, email).Scan(&exists)
	return exists, err
}
