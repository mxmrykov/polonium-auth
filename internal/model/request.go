package model

type (
	SignupCheckRequest struct {
		Email string `json:"email"`
	}

	SignupConfirmCodeRequest struct {
		Email string `json:"email"`
		Pwd   string `json:"pwd"`
		Code  string `json:"code"`
	}

	GetQRCodeRequest struct {
		Email string `json:"email"`
		Pwd   string `json:"pwd"`
	}
)
