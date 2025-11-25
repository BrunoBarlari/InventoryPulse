package service

import (
	"errors"

	"github.com/brunobarlari/inventorypulse/internal/domain/models"
	"github.com/brunobarlari/inventorypulse/internal/repository"
	"github.com/brunobarlari/inventorypulse/pkg/jwt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthorized       = errors.New("unauthorized")
)

type AuthService interface {
	Login(email, password string) (*jwt.TokenPair, error)
	Register(email, password string, role models.Role) (*models.User, error)
	RefreshToken(refreshToken string) (*jwt.TokenPair, error)
	ValidateToken(token string) (*jwt.Claims, error)
}

type authService struct {
	userRepo   repository.UserRepository
	jwtService *jwt.JWTService
}

func NewAuthService(userRepo repository.UserRepository, jwtService *jwt.JWTService) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *authService) Login(email, password string) (*jwt.TokenPair, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.CheckPassword(password) {
		return nil, ErrInvalidCredentials
	}

	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (s *authService) Register(email, password string, role models.Role) (*models.User, error) {
	user := &models.User{
		Email: email,
		Role:  role,
	}

	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) RefreshToken(refreshToken string) (*jwt.TokenPair, error) {
	claims, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, ErrUnauthorized
	}

	// Verify user still exists
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, ErrUnauthorized
	}

	// Generate new token pair
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (s *authService) ValidateToken(token string) (*jwt.Claims, error) {
	return s.jwtService.ValidateToken(token)
}

