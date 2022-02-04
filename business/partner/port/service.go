package port

import (
	"time"
)

type (
	PartnerService struct {
		ID          string    `json:"id"`
		Code        string    `json:"code"`
		Name        string    `json:"name"`
		Pic         string    `json:"pic"`
		Address     string    `json:"address"`
		CallbackUrl string    `json:"callback_url"`
		IpWhitelist []string  `json:"ip_whitelist"`
		Status      string    `json:"status"`
		SecretKey   string    `json:"secret_key"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		DeletedAt   time.Time `json:"deleted_at"`
	}
)

// Service is inbound port
type Service interface {
	//CreateData insert new data
	CreateData(partner PartnerService) error

	//ReadData get data by ID
	ReadData(ID string) (PartnerService, error)

	//UpdateData update new data
	UpdateData(partner PartnerService) error

	//DeleteData delete data
	DeleteData(ID string) error

	//ListData get list data
	ListData() ([]PartnerService, error)
}
