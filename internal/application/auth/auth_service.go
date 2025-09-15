package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/sviatilnik/gophermart/internal/domain/auth"
	"github.com/sviatilnik/gophermart/internal/domain/user"
	"math/rand"
	"time"
)

type Service struct {
	userRepo         user.Repository
	refreshTokenRepo auth.RefreshTokenRepository
	tokenGenerator   TokenGenerator
}

func NewAuthService(repo user.Repository, refreshRepo auth.RefreshTokenRepository, tokenGenerator TokenGenerator) *Service {
	return &Service{
		userRepo:         repo,
		refreshTokenRepo: refreshRepo,
		tokenGenerator:   tokenGenerator,
	}
}

func (s *Service) LoginByID(ctx context.Context, userID string) (*TokenResponse, error) {
	usr, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.GenerateTokens(ctx, usr.ID)
}

func (s *Service) LoginByPassword(ctx context.Context, request LoginRequest) (*TokenResponse, error) {
	if request.Login == "" || request.Password == "" {
		return nil, errors.New("login or password is empty")
	}

	usr, err := s.userRepo.FindByLogin(ctx, user.NewLogin(request.Login))
	if err != nil {
		return nil, err
	}

	err = usr.Password.VerifyPassword(request.Password)
	if err != nil {
		return nil, err
	}

	return s.GenerateTokens(ctx, usr.ID)
}

func (s *Service) LoginByRefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	storedRefreshToken, err := s.refreshTokenRepo.Find(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	if expErr := storedRefreshToken.CheckExpiration(); expErr != nil {
		return nil, expErr
	}

	err = s.refreshTokenRepo.Delete(ctx, storedRefreshToken.Token)
	if err != nil {
		return nil, err
	}

	return s.GenerateTokens(ctx, storedRefreshToken.UserID)
}

func (s *Service) GenerateTokens(ctx context.Context, userID string) (*TokenResponse, error) {
	accessToken, expiredAt, err := s.tokenGenerator.GenerateToken(userID, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpiresAt, err := s.tokenGenerator.GenerateToken(userID, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	refreshToken = refreshToken + s.getRandomString(16)
	h := sha256.New()
	h.Write([]byte(refreshToken))
	refreshToken = base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	err = s.refreshTokenRepo.Save(ctx, &auth.RefreshToken{
		Token:     refreshToken,
		UserID:    userID,
		ExpiresAt: refreshExpiresAt,
	})

	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiredAt.Unix(),
	}, nil
}

func (s *Service) getRandomString(ln int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	b := make([]rune, ln)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}
