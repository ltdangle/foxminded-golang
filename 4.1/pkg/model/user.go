package model

import (
	"time"

	"github.com/google/uuid"
)

// User model.
type User struct {
	Uuid      *string    `db:"uuid"`
	Email     *string    `db:"email"`
	Password  *string    `db:"password"`
	FirstName *string    `db:"first_name"`
	LastName  *string    `db:"last_name"`
	Nickname  *string    `db:"nickname"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
 
func NewUser() *User {
	uuidString := uuid.New().String()
	now := time.Now()
	return &User{
		Uuid:      &uuidString,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}
