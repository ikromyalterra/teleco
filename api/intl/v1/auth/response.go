package auth

type UserResponsetLogin struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
