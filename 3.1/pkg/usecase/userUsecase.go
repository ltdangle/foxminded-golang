package usecase

import (
	"errors"
	"usr_mngmnt/pkg/model"
)

type UserUsecaseInterface interface {
	IsAuthenticated(email string, pass string) bool
	Create(req CreateUserRequest) error
	Update(req CreateUserRequest) error
	View(uuid string) *model.User
	ViewUsers(offset int, limit int) []*model.User
}

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,gte=7"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Nickname  string `json:"nickname" validate:"required"`
}

type userUsecase struct {
	repo model.UserRepoInterface
}

func NewUserUsecases(repo model.UserRepoInterface) *userUsecase {
	return &userUsecase{
		repo: repo,
	}
}

func (ucase *userUsecase) IsAuthenticated(email string, pass string) bool {
	user := ucase.repo.FindByEmailPass(email, pass)

	return user != nil
}

func (ucase *userUsecase) Create(req CreateUserRequest) error {
	u := model.NewUser()
	u.Email = req.Email
	u.FirstName = req.FirstName
	u.LastName = req.LastName
	u.Password = req.Password
	u.Nickname = req.Nickname

	err := ucase.repo.SaveOrUpdate(u)
	if err != nil {
		return err
	}
	return nil
}

func (ucase *userUsecase) Update(req CreateUserRequest) error {
	u := ucase.repo.FindByEmailPass(req.Email, req.Password)
	if u == nil {
		return errors.New("user not found")
	}

	u.Email = req.Email
	u.FirstName = req.FirstName
	u.LastName = req.LastName
	u.Password = req.Password
	u.Nickname = req.Nickname

	err := ucase.repo.SaveOrUpdate(u)
	if err != nil {
		return err
	}
	return nil
}

func (usecase *userUsecase) View(uuid string) *model.User {
	return usecase.repo.Find(uuid)
}

func (usecase *userUsecase) ViewUsers(offset int, limit int) []*model.User {
	return usecase.repo.FindAllUsers(offset, limit)
}
