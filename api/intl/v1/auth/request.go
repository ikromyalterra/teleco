package auth

type UserRequestLogin struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=4,max=20,alphanum"`
	KeepLogin bool   `json:"keep_login"`
}

type UserRefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UserRegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4,max=20,alphanum"`
	Fullname string `json:"fullname" validate:"required"`
}
