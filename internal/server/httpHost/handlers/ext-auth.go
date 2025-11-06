package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
	"github.com/mxmrykov/polonium-auth/internal/model"
	"github.com/mxmrykov/polonium-auth/internal/service"
	"github.com/mxmrykov/polonium-auth/internal/vars"
)

type (
	ExtAuth struct {
		auth service.IAuth
	}
)

func NewExtAuth(auth service.IAuth) *ExtAuth {
	return &ExtAuth{
		auth: auth,
	}
}

func (ea *ExtAuth) SignupCheck(c *gin.Context) {
	ctx := c.Request.Context()

	// ---===Get body===---
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Error: "cannot read request body",
		})
		return
	}

	r := new(model.SignupCheckRequest)
	if err := easyjson.Unmarshal(body, r); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Error: "wrong request body",
		})
		return
	}

	// ---===Check that user can sign up===---
	if err := ea.auth.CanConfirmSignup(ctx, r.Email); err != nil {
		if errors.Is(err, vars.ErrUserAlreadyExists) {
			c.JSON(http.StatusBadRequest, model.Response{
				Error: "user already exists",
			})
			return
		}
		if errors.Is(err, vars.ErrUserAlreadyConfirmingSignup) {
			c.JSON(http.StatusBadRequest, model.Response{
				Error: "user already confirming email",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.Response{
			Error: "unexpected error",
		})
		return
	}

	// ---===Send email code===---
	if err := ea.auth.ConfirmEmail(r.Email); err != nil {
		if errors.Is(err, vars.ErrInvalidEmail) {
			c.JSON(http.StatusBadRequest, model.Response{
				Error: "invalid email",
			})
			return
		}

		c.JSON(http.StatusBadRequest, model.Response{
			Error: "unexpected error",
		})
		return
	}

	c.JSON(http.StatusAccepted, model.Response{
		Message: "verification code has been sent",
	})
}

func (ea *ExtAuth) SignupConfirmEmail(c *gin.Context) {

}
