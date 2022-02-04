package port

type (
	IssuerService struct {
		ID               string `json:"id"`
		Code             string `json:"code"`
		Label            string `json:"label"`
		Config           string `json:"config"`
		ThreadNum        int    `json:"thread_num"`
		ThreadTimeout    int    `json:"thread_timeout"`
		QueueWorkerLimit int    `json:"queue_worker_limit"`
	}
)

// Service is inbound port
type Service interface {
	// CreateData insert new data
	CreateData(issuer IssuerService) error

	// ReadData get data by ID
	ReadData(ID string) (IssuerService, error)

	// UpdateData update new data
	UpdateData(issuer IssuerService) error

	// DeleteData delete data
	DeleteData(ID string) error

	// ListData get list data
	ListData() ([]IssuerService, error)
}
