package model

import (
	"context"
	"database/sql"
	"time"
	"usr_mngmnt/pkg/db"

	_ "github.com/go-sql-driver/mysql"
)

// UserRepoInterface interface.
type UserRepoInterface interface {
	SaveOrUpdate(user *User) error
	Find(uuid string) *User
	FindByEmailPass(email string, pass string) *User
	FindAllUsers(offset int, limit int) []*User
}

// User repo sqlc implementation.
type sqlcUserRepo struct {
	db *db.Queries
}

func NewSqlcRepo(mysqlDsn string) (*sqlcUserRepo, error) {
	mysql, err := sql.Open("mysql", mysqlDsn)
	if err != nil {
		panic(err)
	}

	return &sqlcUserRepo{db: db.New(mysql)}, nil
}

func (r *sqlcUserRepo) FindAllUsers(offset int, limit int) []*User {
	sqlcUsers, err := r.db.FindAllUsers(context.Background(), db.FindAllUsersParams{Offset: int32(offset), Limit: int32(limit)})
	if err != nil {
		return []*User{}
	}

	var users []*User
	for _, sqlcUser := range sqlcUsers {
		users = append(users, r.mapSqlcUserToUser(sqlcUser))
	}

	return users
}

func (r *sqlcUserRepo) SaveOrUpdate(user *User) error {
	u := r.Find(user.Uuid)
	if u == nil {
		err := r.db.CreateUser(context.Background(), db.CreateUserParams{
			Uuid:      sql.NullString{String: user.Uuid, Valid: true},
			Email:     sql.NullString{String: user.Email, Valid: true},
			Password:  sql.NullString{String: user.Password, Valid: true},
			FirstName: sql.NullString{String: user.FirstName, Valid: true},
			LastName:  sql.NullString{String: user.LastName, Valid: true},
			Nickname:  sql.NullString{String: user.Nickname, Valid: true},
			CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
			UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			return err
		}
	} else {
		err := r.db.UpdateUser(context.Background(), db.UpdateUserParams{
			Uuid:      sql.NullString{String: u.Uuid, Valid: true},
			Email:     sql.NullString{String: user.Email, Valid: true},
			Password:  sql.NullString{String: user.Password, Valid: true},
			FirstName: sql.NullString{String: user.FirstName, Valid: true},
			LastName:  sql.NullString{String: user.LastName, Valid: true},
			Nickname:  sql.NullString{String: user.Nickname, Valid: true},
			UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *sqlcUserRepo) Find(uuid string) *User {
	u, err := r.db.GetUser(context.Background(), sql.NullString{String: uuid, Valid: true})
	if err != nil {
		return nil
	}

	return r.mapSqlcUserToUser(u)
}

func (r *sqlcUserRepo) FindByEmailPass(email string, pass string) *User {
	u, err := r.db.FindUserByEmailPass(
		context.Background(),
		db.FindUserByEmailPassParams{
			Email:    sql.NullString{String: email, Valid: true},
			Password: sql.NullString{String: pass, Valid: true},
		})

	if err != nil {
		return nil
	}

	return r.mapSqlcUserToUser(u)
}

func (r *sqlcUserRepo) mapSqlcUserToUser(u db.User) *User {
	return &User{
		Uuid:      u.Uuid.String,
		Email:     u.Email.String,
		Password:  u.Password.String,
		FirstName: u.FirstName.String,
		LastName:  u.LastName.String,
		Nickname:  u.Nickname.String,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
	}

}
