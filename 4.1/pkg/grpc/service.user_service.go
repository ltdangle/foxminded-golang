package grpc

import (
	"context"
	"errors"
	"grpc4_1/pkg/model"
	"log"
	"time"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// UserService is our implementation of the UserServiceServer interface.
type UserService struct {
	UnimplementedUserServiceServer
	repo model.UserRepoInterface
}

func NewUserService(repo model.UserRepoInterface) *UserService {
	return &UserService{repo: repo}
}
func (s *UserService) CreateUser(_ context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	log.Printf("Received CreateUser request: %+v", req)

	// Validate request.
	if err := s.validateVarchar(req.Email); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Email "+err.Error())
	}
	if err := s.validateVarchar(req.Password); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Password "+err.Error())
	}
	if err := s.validateVarchar(req.Firstname); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Firstname "+err.Error())
	}
	if err := s.validateVarchar(req.Lastname); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Lastname "+err.Error())
	}
	if err := s.validateVarchar(req.Nickname); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Nickname "+err.Error())
	}

	user := model.NewUser()
	user.Email = &req.Email
	user.Password = &req.Password
	user.FirstName = &req.Firstname
	user.LastName = &req.Lastname
	user.Nickname = &req.Nickname

	err := s.repo.Create(user)
	if err != nil {
		log.Printf("Create user error: %+v", err)
		return nil, status.Errorf(codes.Internal, "error creating user")
	}

	return &CreateUserResponse{Uuid: *user.Uuid}, nil
}

func (s *UserService) GetUserByEmail(_ context.Context, req *GetUserByEmailRequest) (*GetUserByEmailResponse, error) {
	log.Printf("Received GetUserByEmail request: %+v", req)

	// Validate request.
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "Email field is required")
	}

	user, err := s.repo.FindByEmail(req.GetEmail())
	if err != nil {
		log.Printf("GetUserByEmail error: %+v", err)
		return nil, status.Errorf(codes.Internal, "error retrieving user")
	}

	return &GetUserByEmailResponse{
		User: &User{
			Uuid:      *user.Uuid,
			Email:     *user.Email,
			Firstname: *user.FirstName,
			Lastname:  *user.LastName,
			Nickname:  *user.Nickname,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *UserService) GetUserByID(_ context.Context, req *GetUserByIDRequest) (*GetUserByIDResponse, error) {
	log.Printf("Receive GetUserByID request: %+v", req)

	// Validate request.
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user id field is required")
	}

	user, err := s.repo.FindByUuid(req.GetId())
	if err != nil {
		log.Printf("GetUserById error: %+v", err)
		return nil, status.Errorf(codes.Internal, "error retrieving user")
	}

	return &GetUserByIDResponse{
		User: &User{
			Uuid:      *user.Uuid,
			Email:     *user.Email,
			Firstname: *user.FirstName,
			Lastname:  *user.LastName,
			Nickname:  *user.Nickname,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *UserService) GetUsers(_ context.Context, req *GetUsersRequest) (*GetUsersResponse, error) {
	log.Printf("Receive GetUsers request: %+v", req)
	foundUsers, err := s.repo.FindAll()
	if err != nil {
		log.Printf("GetUsers error: %+v", err)
		return nil, status.Errorf(codes.Internal, "error retrieving users")
	}
	log.Printf("GetUsers result: %+v", foundUsers)

	var userResponses []*User
	for _, foundUser := range foundUsers {
		userResponse := &User{
			Uuid:      *foundUser.Uuid,
			Email:     *foundUser.Email,
			Firstname: *foundUser.FirstName,
			Lastname:  *foundUser.LastName,
			Nickname:  *foundUser.Nickname,
		}
		userResponses = append(userResponses, userResponse)
	}
	return &GetUsersResponse{Users: userResponses}, nil

}

func (s *UserService) UpdateUser(_ context.Context, req *UpdateUserRequest) (*UpdateUserResponse, error) {
	log.Printf("Received CreateUser request: %+v", req)

	// Validate request.
	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "Uuid field is required")
	}
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "Email field is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "Password field is required")
	}
	if req.Firstname == "" {
		return nil, status.Error(codes.InvalidArgument, "Firstname field is required")
	}
	if req.Lastname == "" {
		return nil, status.Error(codes.InvalidArgument, "Lastname field is required")
	}
	if req.Nickname == "" {
		return nil, status.Error(codes.InvalidArgument, "Nickname field is required")
	}

	user := &model.User{}
	user.Uuid = &req.Uuid
	user.Email = &req.Email
	user.Password = &req.Password
	user.FirstName = &req.Firstname
	user.LastName = &req.Lastname
	user.Nickname = &req.Nickname
	now := time.Now()
	user.UpdatedAt = &now

	err := s.repo.Update(user)
	if err != nil {
		log.Printf("Update user error: %+v", err)
		return nil, status.Errorf(codes.Internal, "error updating user")
	}

	return &UpdateUserResponse{Ok: true}, nil
}

func (s *UserService) DeleteUser(_ context.Context, req *DeleteUserRequest) (*DeleteUserResponse, error) {
	log.Printf("Receive DeleteUser request: %+v", req)

	// Validate request.
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Id field is required")
	}

	err := s.repo.Delete(req.Id)
	if err != nil {
		log.Printf("delete user error: %+v", err)
		return nil, status.Errorf(codes.Internal, "error deleting user")
	}

	return &DeleteUserResponse{Ok: true}, nil
}

func (s *UserService) validateVarchar(str string) error {
	if str == "" {
		return errors.New("value must not be empty")
	}
	if len(str) > 255 {
		return errors.New("value must be less than 255 symbols")
	}
	return nil
}
