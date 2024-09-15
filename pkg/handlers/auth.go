package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MobasirSarkar/BeTask/pkg/models"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
)

func (h *Handler) HandleProviderLogin(w http.ResponseWriter, r *http.Request) {
	if u, err := gothic.CompleteUserAuth(w, r); err == nil {
		log.Printf("User already authenticated! %v", u)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func (h *Handler) HandleAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	err = h.auth.StoreUserSession(w, r, user)
	if err != nil {
		log.Println(err)
		return
	}

	userId := generateUUID()

	dbUser := models.User{
		Id:            userId,
		GoogleId:      user.UserID,
		ProfilePicUrl: user.AvatarURL,
		Name:          user.Name,
		Email:         user.Email,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := h.SaveUser(dbUser); err != nil {
		log.Println(err)
		http.Error(w, "Failed to save user data", http.StatusInternalServerError)
		return
	}

	err = h.auth.StoreUserSession(w, r, user)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	log.Println("Logging out...")

	err := gothic.Logout(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	h.auth.RemoveUserSession(w, r)

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) SaveUser(user models.User) error {
	_, err := h.db.Getpool().Exec(
		context.Background(), `
      INSERT INTO users (id, google_id, profile_pic_url, name, email, created_at, updated_at)
      VALUES($1,$2,$3,$4,$5,$6,$7)
      `, user.Id, user.GoogleId, user.ProfilePicUrl, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	return err
}

func generateUUID() string {
	return uuid.NewString()
}
