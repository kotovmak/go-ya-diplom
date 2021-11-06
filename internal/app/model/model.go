package model

type LoginRequest struct {
	Login    string `json:"login" validate:"required,alphanum,max=32"`
	Password string `json:"password" validete:"required,min=8,max=72"`
}
