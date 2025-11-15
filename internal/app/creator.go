package app

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mxmrykov/polonium-auth/internal/auth"
	"github.com/mxmrykov/polonium-auth/internal/repository"
	"github.com/mxmrykov/polonium-auth/internal/server/httpHost/handlers"
	"github.com/mxmrykov/polonium-auth/internal/server/httpHost/middlewares"
	"github.com/mxmrykov/polonium-auth/internal/service"
	"github.com/mxmrykov/polonium-auth/pkg/utils"
)

func (a *Application) setupRoutesAPIV1(
	repositories *repositories,
	jProcessor *auth.JWTProcessor,
) {
	apiV1 := a.httpServer.Router().Group("/ext-auth/api/v1")

	// ---===Middlewares, global setup===---
	{
		apiV1.Use(gin.Recovery(), middlewares.LogMW(), middlewares.CorsMW())
		apiV1.OPTIONS("*any", func(c *gin.Context) {
			c.Writer.WriteHeader(http.StatusOK)
		})
		//apiV1.Use(middlewares.AuthMW(jProcessor, repositories.authRdb))
	}

	// ---===Routing===---
	{
		signupGroup := apiV1.Group("/signup")
		authGroup := apiV1.Group("/auth")
		authService, totpService := service.NewAuth(
			repositories.authPg,
			repositories.authRdb,
			repositories.emailer,
			repositories.vault,
			jProcessor,
		), service.NewTOTP(repositories.vault)
		extAuthHandlers := handlers.NewExtAuth(authService, totpService)
		signupGroup.POST("/general/check", extAuthHandlers.SignupCheck)
		signupGroup.POST("/email/check", extAuthHandlers.SignupConfirmEmail)
		signupGroup.POST("/general/qr", extAuthHandlers.GetQRCode)
		signupGroup.POST("/general/verify", extAuthHandlers.Complete)
		authGroup.POST("/validate", extAuthHandlers.Authorize)
		authGroup.POST("/complete", extAuthHandlers.Complete)
	}
}

func (a *Application) initRepositories() (*repositories, error) {
	authPostgresRepo, err := repository.NewAuthPostgres(&a.cfg.Psql)
	if err != nil {
		return nil, err
	}

	authRedisRepo := repository.NewAuthRedis(&a.cfg.Redis)
	emailer := repository.NewEmailer(&a.cfg.Smtp)
	vault, err := repository.NewAuthVault(&a.cfg.Vault)
	if err != nil {
		return nil, err
	}

	return &repositories{
		authPg:  authPostgresRepo,
		authRdb: authRedisRepo,
		emailer: emailer,
		vault:   vault,
	}, nil
}

func (a *Application) initJwtProcessor() (*auth.JWTProcessor, error) {
	rsa, err := utils.GenerateRSAKeys(1100 + rand.Intn(899))
	if err != nil {
		return nil, fmt.Errorf("cannot init RSA keys: %v", err)
	}

	return auth.NewRSAJWTProcessor(
		rsa.PrivateKey, rsa.PublicKey,
		a.cfg.Auth.Access,
		a.cfg.Auth.Refresh,
		a.cfg.Auth.Issuer,
	), nil
}
