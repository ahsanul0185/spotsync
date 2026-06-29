package user

import (
	"fmt"
	"spotsync/internal/auth"
	"spotsync/internal/domain/user/dto"
)

var ErrInvalidCredentials = fmt.Errorf("invalid email or password")
var ErrUserNotFound = fmt.Errorf("user not found")

type Service interface {
	Register(req dto.RegisterRequest) (*dto.UserResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	GetUserByID(id uint) (*dto.UserResponse, error)
}

type service struct {
	repo       Repository
	jwtService auth.JWTService
}

func NewService(repo Repository, jwtService auth.JWTService) Service {
	return &service{repo: repo, jwtService: jwtService}
}

func (s *service) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {
	user := User{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}
	if user.Role == "" {
		user.Role = "driver"
	}

	if err := user.hashPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.repo.CreateUser(&user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *service) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := user.checkPassword(req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Name, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (s *service) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}
