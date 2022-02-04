package port

import (
	"time"
)

type (
	PartnerRepo struct {
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

// Repository is outbound port
type Repository interface {
	//FindByCode find issuer by code
	FindByCode(code string) PartnerRepo

	//CreateData insert new data
	CreateData(partner PartnerRepo) error

	//ReadData get data by ID
	ReadData(ID string) (PartnerRepo, error)

	//UpdateData update new data
	UpdateData(partner PartnerRepo) error

	//DeleteData delete data
	DeleteData(ID string) error

	//ListData get list data
	ListData() ([]PartnerRepo, error)
}
