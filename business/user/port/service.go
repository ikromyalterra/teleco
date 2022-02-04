package port

type (
	UserService struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Fullname string `json:"fullname"`
		Password string `json:"password"`
	}
)

// Service is inbound port
type Service interface {
	// ReadData get data by ID
	ReadData(ID string) (UserService, error)

	// CreateData insert new data
	CreateData(user UserService) error

	// UpdateData update data
	UpdateData(user UserService) error

	// DeleteData delete data
	DeleteData(ID string) error

	// ListData get list data
	ListData() ([]UserService, error)
}
