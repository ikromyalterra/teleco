package port

type (
	UserTokenRepo struct {
		UserID  string `json:"user_id"`
		TokenID string `json:"token_id"`
	}
)

// Repository is outbound port
type Repository interface {
	// FindByCode find issuer by code
	FindByTokenID(tokenID string) UserTokenRepo

	// CreateData insert new data
	CreateData(userToken UserTokenRepo) error

	// DeleteData revoke user token
	DeleteData(tokenID string) error

	// DeleteDataByUserID remove token by userID
	DeleteDataByUserID(userID string) error
}
