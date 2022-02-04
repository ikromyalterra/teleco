package partner

type ResponsePartner struct {
	ID          string   `json:"id"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Pic         string   `json:"pic"`
	Address     string   `json:"address"`
	CallbackUrl string   `json:"callback_url"`
	IpWhitelist []string `json:"ip_whitelist"`
	Status      string   `json:"status"`
	SecretKey   string   `json:"secret_key"`
}
