package mock

import (
	authPort "github.com/sepulsa/teleco/business/auth/port"

	"github.com/stretchr/testify/mock"
)

type service struct {
	mock.Mock
}

func New() *service {
	return &service{}
}

func (s *service) VerifyPartnerSignature(authData authPort.Signature) (bool, error) {
	result := s.Called(authData)
	return result.Bool(0), result.Error(1)
}

func (s *service) VerifyUserToken(tokenString string) (interface{}, error) {
	result := s.Called(tokenString)
	return result.Get(0).(interface{}), result.Error(1)
}

func (s *service) UserRegister(userAuthService authPort.UserAuthService) error {
	result := s.Called(userAuthService)
	return result.Error(0)
}

func (s *service) UserLogin(userAuthService *authPort.UserAuthService) error {
	result := s.Called(userAuthService)
	return result.Error(0)
}

func (s *service) UserLogout(tokenID string) error {
	result := s.Called(tokenID)
	return result.Error(0)
}

func (s *service) UserRefreshToken(refreshTokenString string) (authPort.UserAuthService, error) {
	result := s.Called(refreshTokenString)
	return result.Get(0).(authPort.UserAuthService), result.Error(1)
}
