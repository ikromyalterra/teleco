package port

type (
	Signature struct {
		Payload     []byte
		PartnerCode string
		Secret      string
		TimeLimit   int
		Token       string
	}

	UserAuthService struct {
		UserID       string
		Email        string
		Password     string
		Fullname     string
		KeepLogin    bool
		Token        string
		RefreshToken string
	}
)

// Service is inbound port
type Service interface {
	// VerifyPartnerSignature verify partner signature
	VerifyPartnerSignature(signature Signature) (bool, error)

	// VerifyUserToken verify jwt
	VerifyUserToken(tokenString string) (interface{}, error)

	// UserRegister register user
	UserRegister(UserAuthService UserAuthService) error

	// UserLogin validate user and generate jwt
	UserLogin(UserAuthService *UserAuthService) error

	// UserLogout revoke user jwt
	UserLogout(tokenID string) error

	// UserRefreshToken refresh jwt
	UserRefreshToken(refreshTokenString string) (UserAuthService, error)
}
