package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func healthcheck(logger *zap.Logger, db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		sqlDB, err := db.DB()
		if err != nil {
			logger.Error("healthcheck: get sql.DB failed", zap.Error(err))
			ctx.Status(http.StatusInternalServerError)
			return
		}

		if err := sqlDB.Ping(); err != nil {
			logger.Error("healthcheck: ping DB failed", zap.Error(err))
			ctx.Status(http.StatusInternalServerError)
		}
	}
}
