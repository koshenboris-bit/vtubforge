package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
		"time"

	"backend-service/internal/config"
	"backend-service/internal/domain"
	"backend-service/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	cfg   config.Config
	users *repository.UserRepository
	tokens *repository.RefreshTokenRepository
}

func NewAuthService(cfg config.Config, repos *repository.Repositories) *AuthService {
	return &AuthService{
		cfg:    cfg,
		users:  repos.Users,
		tokens: repos.Tokens,
	}
}

func (s *AuthService) Register(ctx context.Context, req domain.RegisterRequest) (domain.AuthResponse, error) {
	if req.Login == "" || req.Password == "" {
		return domain.AuthResponse{}, domain.ErrValidation
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.AuthResponse{}, err
	}

	u, err := s.users.Create(ctx, req.Login, string(hash), domain.RoleUser )
	if err != nil {
		return domain.AuthResponse{}, err
	}

	pair, err := s.issueTokens(ctx, u.ID, u.Role)
	if err != nil {
		return domain.AuthResponse{}, err
	}
	return domain.AuthResponse{AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken, User: u}, nil
}

func (s *AuthService) Login(ctx context.Context, req domain.LoginRequest) (domain.AuthResponse, error) {
	u, passwordHash, err := s.users.GetByLogin(ctx, req.Login)
	if err != nil {
		return domain.AuthResponse{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)) != nil {
		return domain.AuthResponse{}, domain.ErrUnauthorized
	}

	pair, err := s.issueTokens(ctx, u.ID, u.Role)
	if err != nil {
		return domain.AuthResponse{}, err
	}
	return domain.AuthResponse{AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken, User: u}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (domain.TokenPair, error) {
	if refreshToken == "" {
		return domain.TokenPair{}, domain.ErrUnauthorized
	}
	userID, err := s.tokens.FindActiveByToken(ctx, refreshToken)
	if err != nil {
		return domain.TokenPair{}, domain.ErrUnauthorized
	}

	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return domain.TokenPair{}, err
	}

	_ = s.tokens.Revoke(ctx, refreshToken)
	return s.issueTokens(ctx, u.ID, u.Role)
}

func (s *AuthService) issueTokens(ctx context.Context, userID int64, role domain.Role) (domain.TokenPair, error) {
	access, err := s.createAccessToken(userID, role)
	if err != nil {
		return domain.TokenPair{}, err
	}
	refresh, err := s.createRefreshToken()
	if err != nil {
		return domain.TokenPair{}, err
	}
	if err := s.tokens.Save(ctx, userID, refresh, time.Now().Add(30*24*time.Hour)); err != nil {
		return domain.TokenPair{}, err
	}
	return domain.TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (s *AuthService) createAccessToken(userID int64, role domain.Role) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": string(role),
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
		"iat":  time.Now().Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(s.cfg.JWTAccessSecret))
}

func (s *AuthService) createRefreshToken() (string, error) {
	b := make([]byte, 48)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func ParseAccessToken(secret, token string) (int64, string, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil || !parsed.Valid {
		return 0, "", domain.ErrInvalidToken
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", domain.ErrInvalidToken
	}

	sub, ok := claims["sub"]
	if !ok {
		return 0, "", domain.ErrInvalidToken
	}

	var userID int64
	switch v := sub.(type) {
	case float64:
		userID = int64(v)
	case int64:
		userID = v
	default:
		return 0, "", domain.ErrInvalidToken
	}

	role, _ := claims["role"].(string)
	if role == "" {
		return 0, "", domain.ErrInvalidToken
	}
	return userID, role, nil
}
