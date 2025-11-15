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
	"github.com/rs/zerolog/log"
)

type (
	ExtAuth struct {
		auth service.IAuth
		totp service.ITOTP
	}
)

func NewExtAuth(
	auth service.IAuth,
	totp service.ITOTP,
) *ExtAuth {
	return &ExtAuth{
		auth: auth,
		totp: totp,
	}
}

func (ea *ExtAuth) SignupCheck(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.Log().Str("logID", c.GetString("logID"))

	// ---===Get body===---
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Err(err).Msg("cannot read body")
		c.JSON(http.StatusBadRequest, model.Response{
			Error: "cannot read request body",
		})
		return
	}

	r := new(model.SignupCheckRequest)
	if err := easyjson.Unmarshal(body, r); err != nil {
		logger.Err(err).Msg("cannot unmarshal request")
		c.JSON(http.StatusBadRequest, model.Response{
			Error: "wrong request body",
		})
		return
	}

	// ---===Check that user can sign up===---
	if err := ea.auth.CanConfirmSignup(ctx, r.Email); err != nil {
		logger.Err(err).Msg("cannot approve confirmation")

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
		logger.Err(err).Msg("cannot confirm email")

		if errors.Is(err, vars.ErrInvalidEmail) {
			c.JSON(http.StatusBadRequest, model.Response{
				Error: "invalid email",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.Response{
			Error: "unexpected error",
		})
		return
	}

	c.JSON(http.StatusAccepted, model.Response{
		Message: "verification code has been sent",
	})
}

func (ea *ExtAuth) SignupConfirmEmail(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.Log().Str("logID", c.GetString("logID"))

	// ---===Get body===---
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Err(err).Msg("cannot read body")
		c.JSON(http.StatusBadRequest, model.Response{
			Error: "cannot read request body",
		})
		return
	}

	r := new(model.SignupConfirmCodeRequest)
	if err := easyjson.Unmarshal(body, r); err != nil {
		logger.Err(err).Msg("cannot unmarshal request")
		c.JSON(http.StatusBadRequest, model.Response{
			Error: "wrong request body",
		})
		return
	}

	// ---===Proceed email verification===---
	if err := ea.auth.ConfirmCode(r.Email, r.Code); err != nil {
		logger.Err(err).Msg("cannot confirm user")

		if errors.Is(err, vars.ErrUserIsNotAuthing) || errors.Is(err, vars.ErrInvalidAuthCode) {
			c.JSON(http.StatusBadRequest, model.Response{
				Error: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.Response{
			Error: "unexpected error",
		})
		return
	}

	// ---===Signup user===---
	if err := ea.auth.SignupUnverified(ctx, r.Email, r.Pwd); err != nil {
		logger.Err(err).Msg("cannot signup user")

		c.JSON(http.StatusInternalServerError, model.Response{
			Error: "unexpected error",
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Message: "verification code processed",
	})
}

func (ea *ExtAuth) GetQRCode(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.Log().Str("logID", c.GetString("logID"))

	// ---===Get body===---
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Err(err).Msg("cannot read body")
		c.JSON(http.StatusBadRequest, model.Response{
			Error: "cannot read request body",
		})
		return
	}

	r := new(model.GetQRCodeRequest)
	if err := easyjson.Unmarshal(body, r); err != nil {
		logger.Err(err).Msg("cannot unmarshal request")
		c.JSON(http.StatusBadRequest, model.Response{
			Error: "wrong request body",
		})
		return
	}

	// ---===Verify that user is valid===---
	if err := ea.auth.VerifyUser(ctx, r.Email, r.Pwd); err != nil {
		logger.Err(err).Msg("cannot verify user")

		if errors.Is(err, vars.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, model.Response{
				Error: "user not found",
			})
			return
		}

		if errors.Is(err, vars.ErrIncorrectPwd) {
			c.JSON(http.StatusUnauthorized, model.Response{
				Error: "invalid password",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.Response{
			Error: "unexpected error",
		})
		return
	}

	// ---===Get QR===---
	QR, err := ea.totp.CreateUserQR(ctx, r.Email)
	if err != nil {
		logger.Err(err).Msg("cannot create QR")
		c.JSON(http.StatusInternalServerError, model.Response{
			Error: "unexpected error",
		})
		return
	}

	c.Data(http.StatusOK, "image/png", QR)
}

func (ea *ExtAuth) Complete(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.Log().Str("logID", c.GetString("logID"))

	// ---===Get body===---
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Err(err).Msg("cannot read body")
		c.JSON(http.StatusBadRequest, model.Response{
			Error: "cannot read request body",
		})
		return
	}
	r := new(model.SignupConfirmCodeRequest)
	if err := easyjson.Unmarshal(body, r); err != nil {
		logger.Err(err).Msg("cannot unmarshal request")
		c.JSON(http.StatusBadRequest, model.Response{
			Error: "wrong request body",
		})
		return
	}

	if err := ea.auth.VerifyUser(ctx, r.Email, r.Pwd); err != nil {
		logger.Err(err).Msg("cannot verify user")

		if errors.Is(err, vars.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, model.Response{
				Error: "user not found",
			})
			return
		}

		if errors.Is(err, vars.ErrIncorrectPwd) {
			c.JSON(http.StatusUnauthorized, model.Response{
				Error: "invalid password",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.Response{
			Error: "unexpected error",
		})
		return
	}

	codeCorrect, err := ea.totp.IsCodeCorrect(ctx, r.Email, r.Code)
	if err != nil {
		logger.Err(err).Msg("cannot verify 2FA code")
		c.JSON(http.StatusInternalServerError, model.Response{
			Error: "unexpected error",
		})
		return
	}

	if !codeCorrect {
		logger.Msg("code is incorrect")
		c.JSON(http.StatusUnauthorized, model.Response{
			Error: "incorrect TOTP code",
		})
		return
	}

	access, refresh, err := ea.auth.CreateSession(r.Email)

	if err != nil {
		logger.Err(err).Msg("cannot create session")
		c.JSON(http.StatusInternalServerError, model.Response{
			Error: "cannot create session",
		})
		return
	}

	if err := ea.auth.VerificateUser(ctx, r.Email); err != nil {
		logger.Err(err).Msg("cannot verificate user")
		c.JSON(http.StatusInternalServerError, model.Response{
			Error: "cannot verificate user",
		})
		return
	}

	c.SetCookie(
		vars.CookiePoloniumAuth,
		refresh,
		3600*24,
		"/",
		"localhost",
		false,
		true,
	)

	c.JSON(http.StatusOK, model.Response{
		Data:    access,
		Message: "user created",
	})
}

func (ea *ExtAuth) Authorize(c *gin.Context) {
	ctx := c.Request.Context()
	logger := log.Log().Str("logID", c.GetString("logID"))

	// ---===Get body===---
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Err(err).Msg("cannot read body")
		c.JSON(http.StatusBadRequest, model.Response{
			Error: "cannot read request body",
		})
		return
	}
	r := new(model.GetQRCodeRequest)
	if err := easyjson.Unmarshal(body, r); err != nil {
		logger.Err(err).Msg("cannot unmarshal request")
		c.JSON(http.StatusBadRequest, model.Response{
			Error: "wrong request body",
		})
		return
	}

	if err := ea.auth.VerifyUser(ctx, r.Email, r.Pwd); err != nil {
		logger.Err(err).Msg("cannot verify user")

		if errors.Is(err, vars.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, model.Response{
				Error: "user not found",
			})
			return
		}

		if errors.Is(err, vars.ErrIncorrectPwd) {
			c.JSON(http.StatusUnauthorized, model.Response{
				Error: "invalid password",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.Response{
			Error: "unexpected error",
		})
		return
	}

	c.JSON(http.StatusAccepted, model.Response{
		Message: "processed",
	})
}
