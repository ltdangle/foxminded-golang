package usecase

import (
	"jwt/pkg/model"
	"time"
)

type UserRequest struct {
	Uuid      string `validate:"required"`
	Email     string `validate:"required"`
	Password  string `validate:"required,gte=7"`
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Nickname  string `validate:"required"`
	RoleId    int    `validate:"required"`
}

// UserUsecaseInterface
type UserUsecaseInterface interface {
	View(uuid string) (model.UserVotes, error)
	Update(*UserRequest) error
	Delete(uuid string) error
	Auth(email string, pass string) (model.UserVotes, error)
	Vote(uuid string, vote int) error
}

// User usecase.
type User struct {
	repo *model.UserRepo
}

func NewUserUsecase(repo *model.UserRepo) *User {
	return &User{repo: repo}
}

func (ucase *User) View(uuid string) (model.UserVotes, error) {
	user, err := ucase.repo.FindByUuid(uuid)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (ucase *User) Update(r *UserRequest) error {
	_, err := ucase.repo.FindByEmailAndPass(r.Email, r.Password)
	if err != nil {
		return err
	}

	user := UserRequestToWriteModel(r)
	err = ucase.repo.Update(user)
	if err != nil {
		return err
	}
	return nil
}

func (ucase *User) Delete(uuid string) error {
	_, err := ucase.repo.FindByUuid(uuid)
	if err != nil {
		return err
	}

	err = ucase.repo.Delete(uuid)
	if err != nil {
		return err
	}
	return nil
}

func (ucase *User) Auth(email string, pass string) (model.UserVotes, error) {
	var user model.UserVotes
	user, err := ucase.repo.FindByEmailAndPass(email, pass)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (ucase *User) Vote(uuid string, vote int) error {
	return nil
}

func UserRequestToWriteModel(r *UserRequest) *model.User {
	updatedAt := time.Now()
	return &model.User{
		Uuid:      &r.Uuid,
		Email:     &r.Email,
		Password:  &r.Password,
		FirstName: &r.FirstName,
		LastName:  &r.LastName,
		Nickname:  &r.Nickname,
		RoleId:    &r.RoleId,
		UpdatedAt: &updatedAt,
	}
}
