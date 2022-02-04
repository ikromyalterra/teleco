package user

type RequestUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"omitempty,min=4,max=20,alphanum"`
	Fullname string `json:"fullname" validate:"required"`
}
