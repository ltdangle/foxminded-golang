package model

type userRepoStub struct {
	user User
}

func NewUserRepoStub() *userRepoStub {
	email := "email@domain.net"
	password := "somepassword"
	firstname := "Firstname"
	lastname := "Lastname"
	nickname := "Nickname"

	user := NewUser()
	user.Email=&email
	user.Password=&password
	user.FirstName=&firstname
	user.LastName=&lastname
	user.Nickname=&nickname

	return &userRepoStub{
		user: *user,
	}
}
func (repo *userRepoStub) Create(_ *User) error {
	return nil
}

func (repo *userRepoStub) Update(_ *User) error {
	return nil
}

func (repo *userRepoStub) Delete(_ string) error {
	return nil
}

func (repo *userRepoStub) FindAll() ([]User, error) {
	return nil, nil
}
func (repo *userRepoStub) FindByUuid(_ string) (User, error) {
	return repo.user, nil
}
func (repo *userRepoStub) FindByEmail(_ string) (User, error) {
	return repo.user, nil

}
