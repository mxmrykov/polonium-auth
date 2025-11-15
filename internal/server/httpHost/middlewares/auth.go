package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mxmrykov/polonium-auth/internal/auth"
	"github.com/mxmrykov/polonium-auth/internal/model"
	"github.com/mxmrykov/polonium-auth/internal/repository"
	"github.com/mxmrykov/polonium-auth/internal/vars"
	"github.com/mxmrykov/polonium-auth/pkg/utils"
	"github.com/rs/zerolog/log"
)

func AuthMW(jp *auth.JWTProcessor, authRdb repository.IAuthRedis) gin.HandlerFunc {
	return func(context *gin.Context) {
		refresh, err := context.Request.Cookie(vars.CookiePoloniumAuth)
		if err != nil {
			log.Log().Msg("no polonium-auth token provided")
			context.AbortWithStatusJSON(http.StatusUnauthorized, model.Response{
				Error: "No polonium-auth token provided",
			})
			return
		}

		claims, err := jp.TokenVerify(refresh.Value)
		if err != nil {
			log.Log().Msg("cannot verify refresh token")
			if errors.Is(err, jwt.ErrTokenExpired) {
				context.AbortWithStatusJSON(http.StatusUnauthorized, model.Response{
					Error: "Session expired",
				})
				return
			}

			context.AbortWithStatusJSON(http.StatusUnauthorized, model.Response{
				Error: "Invalid token",
			})
			return
		}

		access := context.Request.Header.Get(vars.HeaderAuthorization)
		if _, err = jp.TokenVerify(access); err != nil {
			log.Log().Msg("cannot verify access token")
			if errors.Is(err, jwt.ErrTokenExpired) {
				log.Log().Msg("renewing access token")
				session := utils.NewSession()
				if err := authRdb.NewAuthSession(claims.UserID, session); err != nil {
					context.AbortWithStatusJSON(http.StatusServiceUnavailable, model.Response{
						Error: "Cannot renew token",
					})
					return
				}

				newAccess, err := jp.GenerateAccessToken(claims.UserID, session)
				if err != nil {
					context.AbortWithStatusJSON(http.StatusServiceUnavailable, model.Response{
						Error: "Cannot renew token",
					})
					return
				}

				context.JSON(http.StatusResetContent, model.Response{
					Data:    newAccess,
					Message: "Cannot renew token",
				})
				return
			}

			context.AbortWithStatusJSON(http.StatusUnauthorized, model.Response{
				Error: "Invalid access token",
			})
			return
		}
	}
}
