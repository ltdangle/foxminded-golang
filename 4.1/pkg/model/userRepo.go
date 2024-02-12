package model

import (
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type UserRepoInterface interface {
	Create(user *User) error
	Update(user *User) error
	Delete(uuid string) error
	FindAll() ([]User, error)
	FindByUuid(uuid string) (User, error)
	FindByEmail(email string) (User, error)
}

// User model repository.
type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) Create(user *User) error {
	query := `
INSERT INTO users
    ( uuid, email, password, first_name, last_name, nickname, created_at, updated_at)
VALUES 
    ( :uuid, :email, :password, :first_name, :last_name, :nickname, :created_at, :updated_at)`
	_, err := repo.db.NamedExec(query, user)

	if err != nil {
		return err
	}

	return nil
}

func (repo *UserRepo) Update(user *User) error {
	query := `
UPDATE users SET
    email=:email, password:=password, first_name:=first_name, last_name:=last_name, nickname=:nickname, created_at:=created_at, updated_at=:updated_at
WHERE uuid = :uuid; `
	result, err := repo.db.NamedExec(query, user)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("update query returned 0 affected rows")
	}
	if err != nil {
		return err
	}

	return nil
}

func (repo *UserRepo) Delete(uuid string) error {
	query := `
DELETE FROM 
    users 
WHERE uuid = ?;`
	_, err := repo.db.Exec(query, uuid)

	if err != nil {
		return err
	}

	return nil
}

func (repo *UserRepo) FindAll() ([]User, error) {
	var users []User
	query := `
SELECT 
    uuid, email, password, first_name, last_name, nickname, created_at, updated_at  
FROM 
    users;`

	err := repo.db.Select(&users, query)
	if err != nil {
		return users, err
	}
	return users, nil
}
func (repo *UserRepo) FindByUuid(uuid string) (User, error) {
	var user User
	query := `
SELECT 
    uuid, email, password, first_name, last_name, nickname, created_at, updated_at  
FROM 
    users 
WHERE 
  users.uuid = ?;`
	err := repo.db.Get(&user, query, uuid)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (repo *UserRepo) FindByEmail(email string) (User, error) {
	var user User
	query := `
SELECT 
    uuid, email, password, first_name, last_name, nickname, created_at, updated_at  
FROM 
    users 
WHERE 
    users.email=?;
`
	err := repo.db.Get(&user, query, email)
	if err != nil {
		return user, err
	}
	return user, nil
}
