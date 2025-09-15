package auth

import (
	"context"
	"github.com/sviatilnik/gophermart/internal/domain/events"
	"github.com/sviatilnik/gophermart/internal/domain/user"
)

type RegistrationService struct {
	repo            user.Repository
	loginChecker    user.LoginChecker
	passwordChecker user.PasswordChecker
	eventBus        events.Bus
}

func NewRegistrationService(repo user.Repository, bus events.Bus) *RegistrationService {
	return &RegistrationService{
		repo:            repo,
		loginChecker:    user.NewLoginCheckerService(),
		passwordChecker: user.NewSimplePasswordChecker(),
		eventBus:        bus,
	}
}

func (s *RegistrationService) Register(ctx context.Context, req RegisterRequest) (*user.User, error) {
	err := s.loginChecker.Check(req.Login)
	if err != nil {
		return nil, err
	}

	err = s.passwordChecker.Check(req.Password)
	if err != nil {
		return nil, err
	}

	newUser, err := user.NewUser(req.Login, req.Password)
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(ctx, newUser)
	if err != nil {
		return nil, err
	}

	err = s.eventBus.Publish(&user.Registered{UserID: newUser.ID, Email: req.Login})
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
