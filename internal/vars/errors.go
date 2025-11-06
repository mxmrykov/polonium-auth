package vars

import "errors"

var (
	ErrUserAlreadyConfirmingSignup = errors.New("already confirming this email")
	ErrUserAlreadyExists           = errors.New("user with such email already exists")
	ErrInvalidEmail                = errors.New("invalid email")
)
