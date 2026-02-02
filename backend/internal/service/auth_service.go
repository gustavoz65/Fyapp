package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gustavoz65/Cashing-go/internal/config"
	"github.com/gustavoz65/Cashing-go/internal/model"
	"github.com/gustavoz65/Cashing-go/internal/repository"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotActive      = errors.New("user account is not active")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrPasswordMismatch   = errors.New("current password is incorrect")
	ErrWeakPassword       = errors.New("password does not meet requirements")
)

type JWTClaims struct {
	UserID uuid.UUID      `json:"user_id"`
	Email  string         `json:"email"`
	Role   model.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type AuthService struct {
	userRepo *repository.UserRepository
	config   *config.Config
	logger   *zerolog.Logger
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config, logger *zerolog.Logger) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
		logger:   logger,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *model.RegisterRequest) (*model.User, error) {
	// Hash password
	passwordHash, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Email:             req.Email,
		PasswordHash:      passwordHash,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		PreferredCurrency: "BRL",
		PreferredLanguage: "pt-BR",
		Timezone:          "America/Sao_Paulo",
		Role:              model.UserRoleUser,
		IsActive:          true,
	}

	if req.Phone != "" {
		user.Phone = &req.Phone
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create default settings
	_, err = s.userRepo.CreateDefaultSettings(ctx, user.ID)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to create default settings for user")
	}

	s.logger.Info().
		Str("user_id", user.ID.String()).
		Str("email", user.Email).
		Msg("new user registered")

	return user, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest, ipAddress, userAgent string) (*model.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	if !s.checkPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	// Generate tokens
	accessToken, expiresAt, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, refreshTokenHash, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create session
	session := &model.UserSession{
		UserID:           user.ID,
		RefreshTokenHash: refreshTokenHash,
		IPAddress:        &ipAddress,
		UserAgent:        &userAgent,
		ExpiresAt:        time.Now().Add(time.Duration(s.config.Auth.RefreshTokenDuration) * time.Hour),
	}

	if err := s.userRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	s.logger.Info().
		Str("user_id", user.ID.String()).
		Str("email", user.Email).
		Msg("user logged in")

	return &model.LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// RefreshToken refreshes the access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.RefreshTokenResponse, error) {
	tokenHash := s.hashRefreshToken(refreshToken)

	session, err := s.userRepo.GetSessionByToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) ||
			errors.Is(err, repository.ErrSessionExpired) ||
			errors.Is(err, repository.ErrSessionRevoked) {
			return nil, ErrInvalidToken
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Revoke old session
	_ = s.userRepo.RevokeSession(ctx, session.ID)

	// Generate new tokens
	accessToken, expiresAt, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, newRefreshTokenHash, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create new session
	newSession := &model.UserSession{
		UserID:           user.ID,
		RefreshTokenHash: newRefreshTokenHash,
		IPAddress:        session.IPAddress,
		UserAgent:        session.UserAgent,
		ExpiresAt:        time.Now().Add(time.Duration(s.config.Auth.RefreshTokenDuration) * time.Hour),
	}

	if err := s.userRepo.CreateSession(ctx, newSession); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &model.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// Logout revokes the user's refresh token
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := s.hashRefreshToken(refreshToken)

	session, err := s.userRepo.GetSessionByToken(ctx, tokenHash)
	if err != nil {
		// Ignore errors, just return success
		return nil
	}

	return s.userRepo.RevokeSession(ctx, session.ID)
}

// LogoutAll revokes all sessions for a user
func (s *AuthService) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.RevokeAllUserSessions(ctx, userID)
}

// ChangePassword changes the user's password
func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, req *model.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if !s.checkPassword(req.CurrentPassword, user.PasswordHash) {
		return ErrPasswordMismatch
	}

	newPasswordHash, err := s.hashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.userRepo.UpdatePassword(ctx, userID, newPasswordHash); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Revoke all sessions to force re-login
	_ = s.userRepo.RevokeAllUserSessions(ctx, userID)

	s.logger.Info().
		Str("user_id", userID.String()).
		Msg("user changed password")

	return nil
}

// ValidateAccessToken validates an access token and returns the claims
func (s *AuthService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Auth.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// Helper methods

func (s *AuthService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s *AuthService) checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) generateAccessToken(user *model.User) (string, int64, error) {
	expiresAt := time.Now().Add(time.Duration(s.config.Auth.AccessTokenDuration) * time.Minute)

	claims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.config.Auth.Issuer,
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.Auth.SecretKey))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt.Unix(), nil
}

func (s *AuthService) generateRefreshToken() (string, string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}

	token := base64.URLEncoding.EncodeToString(bytes)
	hash := s.hashRefreshToken(token)

	return token, hash, nil
}

func (s *AuthService) hashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
