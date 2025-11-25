package jwt

import (
	"errors"
	"time"

	"github.com/brunobarlari/inventorypulse/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type JWTService struct {
	secret            string
	expiryHours       int
	refreshExpiryHours int
}

func NewJWTService(cfg *config.JWTConfig) *JWTService {
	return &JWTService{
		secret:            cfg.Secret,
		expiryHours:       cfg.ExpiryHours,
		refreshExpiryHours: cfg.RefreshExpiryHours,
	}
}

// GenerateTokenPair creates both access and refresh tokens
func (s *JWTService) GenerateTokenPair(userID uint, email, role string) (*TokenPair, error) {
	accessToken, expiresAt, err := s.generateToken(userID, email, role, s.expiryHours)
	if err != nil {
		return nil, err
	}

	refreshToken, _, err := s.generateToken(userID, email, role, s.refreshExpiryHours)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// GenerateAccessToken creates a new access token
func (s *JWTService) GenerateAccessToken(userID uint, email, role string) (string, int64, error) {
	return s.generateToken(userID, email, role, s.expiryHours)
}

func (s *JWTService) generateToken(userID uint, email, role string, expiryHours int) (string, int64, error) {
	expirationTime := time.Now().Add(time.Duration(expiryHours) * time.Hour)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "inventorypulse",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expirationTime.Unix(), nil
}

// ValidateToken validates the token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

