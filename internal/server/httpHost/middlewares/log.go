package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func LogMW() gin.HandlerFunc {
	return func(context *gin.Context) {
		log.Log().
			Str("method", context.Request.Method).
			Str("path", context.Request.RequestURI).
			Msg("handle")
	}
}
