package model

import "time"

// User role descriptions.
const ROLE_USER = "user"
const ROLE_ADMIN = "admin"
const ROLE_MODERATOR = "moderator"

// User write model.
type User struct {
	Uuid      *string    `db:"uuid"`
	Email     *string    `db:"email"`
	Password  *string    `db:"password"`
	FirstName *string    `db:"first_name"`
	LastName  *string    `db:"last_name"`
	Nickname  *string    `db:"nickname"`
	RoleId    *int       `db:"role_id"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

// UserVotes read model.
type UserVotes struct {
	Uuid      *string    `db:"uuid"`
	Email     *string    `db:"email"`
	Password  *string    `db:"password"`
	FirstName *string    `db:"first_name"`
	LastName  *string    `db:"last_name"`
	Nickname  *string    `db:"nickname"`
	Role      *string    `db:"role"`
	Votes     *string    `db:"votes"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

// Vote read model.
type Vote struct {
	FromUser  *string    `db:"from_user"`
	ToUser    *string    `db:"to_user"`
	Vote      *int       `db:"vote"`
	UpdatedAt *time.Time `db:"updated_at"`
}
