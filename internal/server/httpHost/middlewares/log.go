package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func LogMW() gin.HandlerFunc {
	return func(context *gin.Context) {
		t := time.Now()
		logID := uuid.New().String()
		context.Set("logID", logID)
		context.Next()
		log.Log().
			Str("method", context.Request.Method).
			Str("path", context.Request.RequestURI).
			Str("duration", time.Since(t).String()).
			Str("logID", logID).
			Msg("handle")
	}
}
