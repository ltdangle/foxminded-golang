package usecase

import "usr_mngmnt/pkg/model"

type UserInterfaceStub struct {
	viewFn func() *model.User
}

func NewUserUsecaseInterfaceStub() *UserInterfaceStub {
	return &UserInterfaceStub{}
}

func (stub *UserInterfaceStub) IsAuthenticated(_ string, _ string) bool {
	return false
}

func (stub *UserInterfaceStub) Create(_ CreateUserRequest) error {
	return nil
}

func (stub *UserInterfaceStub) Update(_ CreateUserRequest) error {
	return nil
}

func (stub *UserInterfaceStub) SetView(f func() *model.User) {
	stub.viewFn = f
}

func (stub *UserInterfaceStub) View(_ string) *model.User {
	return stub.viewFn()
}

func (stub *UserInterfaceStub) ViewUsers(_ int, _ int) []*model.User {
	return nil
}
