package services

import (
	"biletter-service/internal/models"
	"biletter-service/internal/repository"
	"fmt"
)

type UserService interface {
	GetByEmail(email string) (*models.User, error)
	GetByID(userID int) (*models.User, error)
	ValidateCredentials(email, password string) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (s *userService) GetByID(userID int) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return user, nil
}

func (s *userService) ValidateCredentials(email, password string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("user is not active")
	}

	// В Java используется plaintext password, следую той же логике
	if user.PasswordPlain == nil || *user.PasswordPlain != password {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}
