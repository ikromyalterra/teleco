package auth_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strconv"
	"testing"
	"time"

	mockPartnerRepo "github.com/sepulsa/teleco/modules/repository/mock/partner"
	mockUserRepo "github.com/sepulsa/teleco/modules/repository/mock/user"
	mockUserTokenRepo "github.com/sepulsa/teleco/modules/repository/mock/usertoken"

	authService "github.com/sepulsa/teleco/business/auth"
	authPort "github.com/sepulsa/teleco/business/auth/port"
	partnerPort "github.com/sepulsa/teleco/business/partner/port"
	userPort "github.com/sepulsa/teleco/business/user/port"
	"github.com/sepulsa/teleco/utils/auth"
	"github.com/sepulsa/teleco/utils/minifier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	TestID             = "6138813fb95630b0b528b160"
	TestEmail          = "test@test.com"
	TestPassword       = "ikromy"
	TestPasswordStored = "$2a$08$b9Ewo/XS5/lPlviA6JwzheYx.zi9erD3.x9ip6ZJDjPfB.SSqEo5G"
	TestEmailDuplicate = "duplicate@test.com"
	TestTokenID        = "6328402206636078"
	TestBadToken       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	TestErrorDuplicateEmail       error = authService.ErrUserEmailDuplicate
	TestErrorBadCredential        error = authService.ErrInvalidCredential
	TestErrorInvalidJWTCredential error = errors.New("invalid or expired token")

	TestPartnerCode      = "partnercode"
	TestPartnerSecretKey = "partner-secret-key"
)

func TestUserRegister(t *testing.T) {
	userRepo := mockUserRepo.New()
	userTokenRepo := mockUserTokenRepo.New()

	// start populate mock data
	outputUserFound := userPort.UserRepo{
		ID:    TestID,
		Email: TestEmail,
	}
	outputUserNotFound := userPort.UserRepo{
		ID: "",
	}
	inputUserServDuplicateEmail := authPort.UserAuthService{
		Email: TestEmailDuplicate,
	}
	inputUserServ := authPort.UserAuthService{
		Email:    TestEmail,
		Password: TestPassword,
	}
	userRepo.On("FindByEmail", TestEmailDuplicate).Return(outputUserFound)
	userRepo.On("FindByEmail", TestEmail).Return(outputUserNotFound)
	userRepo.On("CreateData", mock.Anything).Return(nil)

	s := authService.New(userRepo, userTokenRepo, nil)

	// success
	assert.Nil(t, s.UserRegister(inputUserServ))

	// err duplicate email
	err := s.UserRegister(inputUserServDuplicateEmail)
	assert.Equal(t, TestErrorDuplicateEmail, err)
}

func TestUserLogin(t *testing.T) {
	userRepo := mockUserRepo.New()
	userTokenRepo := mockUserTokenRepo.New()

	dataUserRepo := userPort.UserRepo{
		ID:       TestID,
		Email:    TestEmail,
		Password: TestPasswordStored,
	}

	userRepo.On("FindByEmail", mock.Anything).Return(dataUserRepo).Once()
	userRepo.On("FindByEmail", mock.Anything).Return(userPort.UserRepo{ID: ""}).Once()
	userTokenRepo.On("CreateData", mock.Anything).Return(nil)

	s := authService.New(userRepo, userTokenRepo, nil)

	inputUserLogin := new(authPort.UserAuthService)
	inputUserLogin.Email = TestEmail
	inputUserLogin.Password = TestPassword
	inputUserLogin.KeepLogin = true

	// success
	assert.Nil(t, s.UserLogin(inputUserLogin))

	// err bad credential email
	err := s.UserLogin(inputUserLogin)
	assert.Equal(t, TestErrorBadCredential, err)
}

func TestUserLogout(t *testing.T) {
	userRepo := mockUserRepo.New()
	userTokenRepo := mockUserTokenRepo.New()

	userTokenRepo.On("DeleteData", mock.Anything).Return(nil)

	s := authService.New(userRepo, userTokenRepo, nil)

	assert.Nil(t, s.UserLogout(TestTokenID))
}

func TestUserRefreshToken(t *testing.T) {
	userRepo := mockUserRepo.New()
	userTokenRepo := mockUserTokenRepo.New()

	// populate data
	dataUserTokenRepo := authPort.UserTokenRepo{
		UserID:  TestID,
		TokenID: TestTokenID,
	}
	dataUserRepo := userPort.UserRepo{
		ID:    TestID,
		Email: TestEmail,
	}
	userTokenRepo.On("FindByTokenID", mock.Anything).Return(authPort.UserTokenRepo{}).Once()
	userTokenRepo.On("FindByTokenID", mock.Anything).Return(dataUserTokenRepo)
	userTokenRepo.On("CreateData", mock.Anything).Return(errors.New("")).Once()
	userTokenRepo.On("CreateData", mock.Anything).Return(nil)
	userTokenRepo.On("DeleteData", mock.Anything).Return(nil)
	userRepo.On("ReadData", mock.Anything).Return(dataUserRepo, nil)

	// inject
	s := authService.New(userRepo, userTokenRepo, nil)

	// err parse token
	_, err := s.UserRefreshToken(TestBadToken)
	assert.Equal(t, TestErrorInvalidJWTCredential, err)

	// err user token not found
	_, err = s.UserRefreshToken(auth.TestJWTAlwaysValid)
	assert.Equal(t, authService.ErrInvalidToken, err)

	// err while update user token
	_, err = s.UserRefreshToken(auth.TestJWTAlwaysValid)
	assert.Equal(t, authService.ErrGenerateToken, err)

	// success
	_, err = s.UserRefreshToken(auth.TestJWTAlwaysValid)
	assert.Nil(t, err)
}

func TestVerifyUserToken(t *testing.T) {
	userRepo := mockUserRepo.New()
	userTokenRepo := mockUserTokenRepo.New()

	dataUserTokenRepo := authPort.UserTokenRepo{
		UserID:  TestID,
		TokenID: TestTokenID,
	}
	userTokenRepo.On("FindByTokenID", mock.Anything).Return(authPort.UserTokenRepo{}).Once()
	userTokenRepo.On("FindByTokenID", mock.Anything).Return(dataUserTokenRepo)

	s := authService.New(userRepo, userTokenRepo, nil)

	// err bad jwt
	_, err := s.VerifyUserToken(TestBadToken)
	assert.Equal(t, TestErrorInvalidJWTCredential, err)

	// err user token not found
	_, err = s.VerifyUserToken(auth.TestJWTAlwaysValid)
	assert.Equal(t, authService.ErrInvalidToken, err)

	// success
	_, err = s.VerifyUserToken(auth.TestJWTAlwaysValid)
	assert.Nil(t, err)
}

func TestVerifyPartnerSignature(t *testing.T) {
	// generate token first
	payload := `{
			"order_id": "123"
		}`
	now := time.Now()
	fiveSecondsAgo := time.Second * time.Duration(-5)
	past := now.Add(fiveSecondsAgo).Unix()
	strPast := strconv.FormatInt(past, 10)
	payloadMinified, _ := minifier.JSON([]byte(payload))
	payloadBeforeHashed := strPast + ":" + string(payloadMinified)
	h := hmac.New(sha256.New, []byte(TestPartnerSecretKey))
	h.Write([]byte(payloadBeforeHashed))
	payloadHashed := hex.EncodeToString(h.Sum(nil))
	dataSignature := strPast + ":" + payloadHashed
	signature := base64.URLEncoding.EncodeToString([]byte(dataSignature))
	// end generate token

	dataService := authPort.Signature{
		Payload:     []byte(payloadMinified),
		Token:       signature,
		PartnerCode: TestPartnerCode,
		TimeLimit:   60,
	}

	partnerRepo := mockPartnerRepo.New()
	partnerRepo.On("FindByCode", mock.Anything).Return(partnerPort.PartnerRepo{}).Once()
	partnerRepo.On("FindByCode", mock.Anything).Return(partnerPort.PartnerRepo{SecretKey: TestPartnerSecretKey}).Once()

	// inject
	service := authService.New(nil, nil, partnerRepo)

	// error partner secret key not
	valid, err := service.VerifyPartnerSignature(dataService)
	assert.NotNil(t, err)
	assert.Equal(t, false, valid)

	// success
	valid, err = service.VerifyPartnerSignature(dataService)
	assert.Nil(t, err)
	assert.Equal(t, true, valid)
}
