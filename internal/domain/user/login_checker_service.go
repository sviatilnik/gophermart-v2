package user

type LoginChecker interface {
	Check(login string) error
}

type LoginCheckerService struct{}

func NewLoginCheckerService() *LoginCheckerService {
	return &LoginCheckerService{}
}

func (s *LoginCheckerService) Check(login string) error {
	if len(login) < 3 || len(login) > 128 {
		return ErrLoginNotValid
	}

	return nil
}
