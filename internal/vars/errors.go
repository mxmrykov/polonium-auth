package vars

import "errors"

var (
	ErrUserAlreadyConfirmingSignup = errors.New("already confirming this email")
	ErrUserAlreadyExists           = errors.New("user with such email already exists")
	ErrInvalidEmail                = errors.New("invalid email")
	ErrUserIsNotAuthing            = errors.New("user isn`t authorizing sessions or code is expired")
	ErrInvalidAuthCode             = errors.New("invalid auth code")
	ErrUserNotFound                = errors.New("user does not exists")
	ErrNoSuchVariableInVault       = errors.New("no such variable in vault")
	ErrIncorrectPwd                = errors.New("incorrect password")
	ErrUserAlreadyVerified         = errors.New("user is already verified")
)
