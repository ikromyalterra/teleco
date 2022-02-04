package auth

import (
	"errors"
	"strings"

	authPort "github.com/sepulsa/teleco/business/auth/port"
	partnerPort "github.com/sepulsa/teleco/business/partner/port"
	userPort "github.com/sepulsa/teleco/business/user/port"
	"github.com/sepulsa/teleco/utils/auth"
	"github.com/sepulsa/teleco/utils/crypto"
	"github.com/sepulsa/teleco/utils/helper"
	"github.com/sepulsa/teleco/utils/validator"
)

type (
	service struct {
		userRepository      userPort.Repository
		userTokenRepository authPort.Repository
		partnerRepository   partnerPort.Repository
	}
)

func New(uRepo userPort.Repository, uTokenRepo authPort.Repository, pRepo partnerPort.Repository) authPort.Service {
	return &service{
		uRepo,
		uTokenRepo,
		pRepo,
	}
}

var (
	jwt                     auth.JWT = auth.NewJWT()
	ErrUserEmailDuplicate   error    = errors.New("email already in use")
	ErrGenerateToken        error    = errors.New("generate token failed")
	ErrInvalidToken         error    = errors.New("invalid token")
	ErrInvalidCredential    error    = errors.New("invalid credential")
	ErrUserGeneratePassword error    = errors.New("generate password failed")
)

func (s *service) VerifyPartnerSignature(authData authPort.Signature) (bool, error) {
	partnerData := s.partnerRepository.FindByCode(authData.PartnerCode)
	if strings.TrimSpace(partnerData.SecretKey) == "" {
		return false, ErrInvalidCredential
	}
	signature := validator.Signature{
		Token:     authData.Token,
		Payload:   authData.Payload,
		TimeLimit: authData.TimeLimit,
		Secret:    partnerData.SecretKey,
	}

	validate := validator.NewSignatureValidator(signature)

	return validate.Bind(validator.Parse).
		Bind(validator.VerifyReqTime).
		Bind(validator.VerifyReqPayload).
		Verify()
}

func (s *service) VerifyUserToken(tokenString string) (interface{}, error) {
	token, tokenClaims, err := jwt.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	userToken := s.userTokenRepository.FindByTokenID(tokenClaims.ID)
	if userToken.UserID == "" {
		return nil, ErrInvalidToken
	}
	return token, nil
}

func (s *service) UserRegister(userAuth authPort.UserAuthService) error {
	existingUserEmail := s.userRepository.FindByEmail(userAuth.Email)
	if existingUserEmail.ID != "" {
		return ErrUserEmailDuplicate
	}

	hashedPassword, err := crypto.UserGeneratePassword(userAuth.Password)
	if err != nil {
		return ErrUserGeneratePassword
	}

	var data userPort.UserRepo
	data.Email = userAuth.Email
	data.Fullname = userAuth.Fullname
	data.Password = hashedPassword

	return s.userRepository.CreateData(data)
}

func (s *service) UserLogin(userAuth *authPort.UserAuthService) (err error) {
	if !s.bindUserCredential(userAuth) {
		return ErrInvalidCredential
	}

	tokenID := helper.GenerateUniqueID()
	err = s.generateToken(tokenID, userAuth)
	if err == nil {
		if userAuth.KeepLogin {
			err = s.generateRefreshToken(tokenID, userAuth)
		}
		if err == nil {
			var userTokenRepo authPort.UserTokenRepo
			userTokenRepo.TokenID = tokenID
			userTokenRepo.UserID = userAuth.UserID

			return s.userTokenRepository.CreateData(userTokenRepo)
		}
	}

	return ErrGenerateToken
}

func (s *service) UserLogout(tokenID string) error {
	return s.userTokenRepository.DeleteData(tokenID)
}

func (s *service) UserRefreshToken(refreshTokenString string) (userAuth authPort.UserAuthService, err error) {
	refreshTokenClaims, err := jwt.ParseRefreshToken(refreshTokenString)
	if err != nil {
		return
	}
	if err = s.populateUserToken(refreshTokenClaims.ID, &userAuth); err == nil {
		newTokenID := helper.GenerateUniqueID()
		if err = s.generateTokens(newTokenID, &userAuth); err == nil {
			err = s.updateUserToken(refreshTokenClaims.ID, newTokenID, userAuth.UserID)
		}
	}

	return
}

func (s *service) generateToken(tokenID string, userAuth *authPort.UserAuthService) (err error) {
	var tokenClaims auth.JWTClaims
	tokenClaims.ID = tokenID
	tokenClaims.Email = userAuth.Email
	tokenClaims.Fullname = userAuth.Fullname

	userAuth.Token, err = jwt.CreateToken(tokenClaims)

	return
}

func (s *service) generateRefreshToken(tokenID string, userAuth *authPort.UserAuthService) (err error) {
	var refreshTokenClaims auth.JWTRefreshClaims
	refreshTokenClaims.ID = tokenID

	userAuth.RefreshToken, err = jwt.CreateRefreshToken(refreshTokenClaims)

	return
}

func (s *service) updateUserToken(oldTokenID, newtokenID, userID string) error {
	var userTokenRepo authPort.UserTokenRepo

	userTokenRepo.TokenID = newtokenID
	userTokenRepo.UserID = userID

	if err := s.userTokenRepository.CreateData(userTokenRepo); err == nil {
		return s.userTokenRepository.DeleteData(oldTokenID)
	}

	return ErrGenerateToken
}

func (s *service) generateTokens(tokenID string, userAuth *authPort.UserAuthService) error {
	if err := s.generateToken(tokenID, userAuth); err == nil {
		return s.generateRefreshToken(tokenID, userAuth)
	}

	return ErrGenerateToken
}

func (s *service) populateUserToken(tokenID string, userAuth *authPort.UserAuthService) error {
	token := s.userTokenRepository.FindByTokenID(tokenID)
	if token.UserID != "" {
		userAuth.UserID = token.UserID
		user, err := s.userRepository.ReadData(userAuth.UserID)
		if err == nil {
			userAuth.Email = user.Email
			userAuth.Fullname = user.Fullname
			return nil
		}
	}

	return ErrInvalidToken
}

func (s *service) bindUserCredential(user *authPort.UserAuthService) bool {
	existingUser := s.userRepository.FindByEmail(user.Email)
	if existingUser.ID != "" {
		user.UserID = existingUser.ID
		user.Fullname = existingUser.Fullname
		return crypto.UserVerifyPassword(user.Password, existingUser.Password)
	}

	return false
}
