package handlers

import "github.com/gin-gonic/gin"

type (
	ExtAuth struct {
	}
)

func NewExtAuth() *ExtAuth {
	return &ExtAuth{}
}

func (ea *ExtAuth) SignupCheck(c *gin.Context) {

}

func (ea *ExtAuth) SignupConfirmEmail(c *gin.Context) {

}
