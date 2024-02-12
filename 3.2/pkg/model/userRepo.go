package model

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func NewUser() *User {
	uuidString := uuid.New().String()
	now := time.Now()
	return &User{
		Uuid:      &uuidString,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}

// User write model repository.
type UserRepo struct {
	db *sqlx.DB
}

func NewWriteUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) Save(user *User) error {
	query := `
INSERT INTO users
    ( uuid, email, password, first_name, last_name, nickname, role_id, created_at, updated_at)
VALUES 
    ( :uuid, :email, :password, :first_name, :last_name, :nickname, :role_id, :created_at, :updated_at)`
	_, err := repo.db.NamedExec(query, user)

	if err != nil {
		return err
	}

	return nil
}

func (repo *UserRepo) Update(user *User) error {
	query := `
UPDATE users SET
    email=:email, password:=password, first_name:=first_name, last_name:=last_name, nickname=:nickname, role_id=:role_id, created_at:=created_at, updated_at=:updated_at
WHERE uuid = :uuid; `
	_, err := repo.db.NamedExec(query, user)

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

func (repo *UserRepo) FindByUuid(uuid string) (UserVotes, error) {
	var user UserVotes
	query := `
SELECT 
    users.uuid, users.email, users.password, users.first_name, users.last_name, users.nickname, users.created_at, users.updated_at, user_role.title AS role, SUM(vote.vote) AS votes
FROM 
    users 
LEFT JOIN 
    user_role ON users.role_id = user_role.id
LEFT JOIN vote ON users.uuid=vote.to_user
WHERE 
    users.uuid = ? 
GROUP BY users.uuid;`
	err := repo.db.Get(&user, query, uuid)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (repo *UserRepo) FindByEmailAndPass(email string, pass string) (UserVotes, error) {
	var user UserVotes
	query := `
SELECT 
    uuid, email, password, first_name, last_name, nickname, created_at, updated_at, user_role.title AS role 
FROM 
    users 
LEFT JOIN 
    user_role ON users.role_id = user_role.id
WHERE 
    users.email=? and users.password=?;
`
	err := repo.db.Get(&user, query, email, pass)
	if err != nil {
		return user, err
	}
	return user, nil
}
