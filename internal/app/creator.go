package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mxmrykov/polonium-auth/internal/repository"
	"github.com/mxmrykov/polonium-auth/internal/server/httpHost/handlers"
	"github.com/mxmrykov/polonium-auth/internal/server/httpHost/middlewares"
	"github.com/mxmrykov/polonium-auth/internal/service"
)

func (a *Application) setupRoutesAPIV1(
	authPg repository.IAuthPostgres,
	authRdb repository.IAuthRedis,
	emailer repository.IEmailer,
) {
	apiV1 := a.httpServer.Router().Group("/ext-auth/api/v1")

	{
		apiV1.Use(gin.Recovery(), middlewares.LogMW(), middlewares.CorsMW())
		apiV1.OPTIONS("*any", func(c *gin.Context) {
			c.Writer.WriteHeader(http.StatusOK)
		})
		apiV1.Use(middlewares.AuthMW())
	}

	{
		authService := service.NewAuth(
			authPg,
			authRdb,
			emailer,
		)
		extAuthHandlers := handlers.NewExtAuth(authService)
		apiV1.POST("/signup/general/check", extAuthHandlers.SignupCheck)
		apiV1.POST("/signup/email/check", extAuthHandlers.SignupConfirmEmail)
	}
}

func (a *Application) initRepositories() (*repositories, error) {
	authPostgresRepo, err := repository.NewAuthPostgres(&a.cfg.Psql)
	if err != nil {
		return nil, err
	}

	authRedisRepo := repository.NewAuthRedis(&a.cfg.Redis)
	emailer := repository.NewEmailer(&a.cfg.Smtp)

	return &repositories{
		authPg:  authPostgresRepo,
		authRdb: authRedisRepo,
		emailer: emailer,
	}, nil
}
