package port

type (
	IssuerRepo struct {
		ID               string `json:"id"`
		Code             string `json:"code"`
		Label            string `json:"label"`
		Config           string `json:"config"`
		ThreadNum        int    `json:"thread_num"`
		ThreadTimeout    int    `json:"thread_timeout"`
		QueueWorkerLimit int    `json:"queue_worker_limit"`
	}
)

// Repository is outbound port
type Repository interface {
	//FindByCode find issuer by code
	FindByCode(code string) IssuerRepo

	//CreateData insert new data
	CreateData(issuer IssuerRepo) error

	//ReadData get data by ID
	ReadData(ID string) (IssuerRepo, error)

	//UpdateData update new data
	UpdateData(issuer IssuerRepo) error

	//DeleteData delete data
	DeleteData(ID string) error

	//ListData get list data
	ListData() ([]IssuerRepo, error)
}
