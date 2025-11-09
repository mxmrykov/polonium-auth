package model

import "time"

type (
	User struct {
		Email, Id, SshSign, Deployer string
		Verified, Banned             bool
		CreateDt                     time.Time
	}
)
