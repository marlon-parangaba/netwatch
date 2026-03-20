package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/ertel/netwatch-backend/internal/config"
	"github.com/ertel/netwatch-backend/internal/models"
	"github.com/ertel/netwatch-backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrUserInactive       = errors.New("usuário inativo")
	ErrTokenInvalid       = errors.New("token inválido ou expirado")
	ErrUnauthorized       = errors.New("sem permissão para esta ação")
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthService interface {
	Login(req LoginRequest) (*TokenPair, *models.User, error)
	RefreshToken(refreshToken string) (*TokenPair, error)
	ValidateToken(token string) (*JWTClaims, error)
	Register(req models.CreateUserRequest) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{userRepo: userRepo, cfg: cfg}
}

func (s *authService) Login(req LoginRequest) (*TokenPair, *models.User, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, nil, ErrInvalidCredentials
		}
		return nil, nil, err
	}

	if !user.Active {
		return nil, nil, ErrUserInactive
	}

	if !user.CheckPassword(req.Password) {
		return nil, nil, ErrInvalidCredentials
	}

	tokens, err := s.generateTokenPair(user)
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao gerar tokens: %w", err)
	}

	_ = s.userRepo.UpdateLastLogin(user.ID)

	return tokens, user, nil
}

func (s *authService) RefreshToken(refreshToken string) (*TokenPair, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, ErrTokenInvalid
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, ErrTokenInvalid
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil || !user.Active {
		return nil, ErrTokenInvalid
	}

	return s.generateTokenPair(user)
}

func (s *authService) ValidateToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", t.Header["alg"])
		}
		return []byte(s.cfg.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

func (s *authService) Register(req models.CreateUserRequest) (*models.User, error) {
	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}

	if err := user.SetPassword(req.Password); err != nil {
		return nil, fmt.Errorf("erro ao gerar hash de senha: %w", err)
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("erro ao criar usuário: %w", err)
	}

	return user, nil
}

func (s *authService) generateTokenPair(user *models.User) (*TokenPair, error) {
	expiration := time.Duration(s.cfg.JWT.ExpirationHours) * time.Hour
	expiresAt := time.Now().Add(expiration)

	claims := &JWTClaims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, err
	}

	refreshExpiration := time.Duration(s.cfg.JWT.RefreshExpirationHours) * time.Hour
	refreshClaims := &JWTClaims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Unix(),
	}, nil
}
