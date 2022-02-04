package port

type (
	UserRepo struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Fullname string `json:"fullname"`
		Password string `json:"password"`
	}
)

// Repository is outbound port
type Repository interface {
	// FindByEmail find issuer by code
	FindByEmail(email string) UserRepo

	// CreateData insert new data
	CreateData(user UserRepo) error

	UpdateData(user UserRepo) error

	// ReadData get data by ID
	ReadData(ID string) (UserRepo, error)

	// DeleteData delete data
	DeleteData(ID string) error

	// ListData get list data
	ListData() ([]UserRepo, error)
}
