package models

import (
	"time"
)

type User struct {
	Id            string    `json:"id"`
	GoogleId      string    `json:"google_id"`
	ProfilePicUrl string    `json:"profile_pic_url"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
