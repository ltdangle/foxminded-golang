package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Uuid      string
	Email     string
	Password  string
	FirstName string
	LastName  string
	Nickname  string 
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser() *User {
	return &User{
		Uuid:      uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
