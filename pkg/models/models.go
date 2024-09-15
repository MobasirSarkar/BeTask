package models

import (
	"time"
)

type User struct {
	Id            string
	GoogleId      string
	ProfilePicUrl string
	Name          string
	Email         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
