package partner

type RequestPartner struct {
	Code        string   `json:"code" validate:"required"`
	Name        string   `json:"name" validate:"required"`
	Pic         string   `json:"pic" validate:"required"`
	Address     string   `json:"address"`
	CallbackUrl string   `json:"callback_url"`
	IpWhitelist []string `json:"ip_whitelist"`
	Status      string   `json:"status"`
	SecretKey   string   `json:"secret_key"`
}
