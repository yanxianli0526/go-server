package router

import (
	"meepShopTest/internal/database"
	"meepShopTest/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func New(router *gin.Engine,
	logger *zap.Logger,
	db *database.GormDatabase) *gin.Engine {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))
	router.GET("/api/readyz", healthcheck(logger, db.DB))

	// 題目沒有特別說要做登入之類的動作就沒有特別做了 不過如果有要登入的話這邊可以塞個檢查身份相關的middleware
	// 或是如果我有誤會意思 其實有需要登入的話可以再和我說 我會再補上跟改寫
	router.Use(middleware.ErrorHandler)

	apiRouter := router.Group("/api")

	RegisterUser(db, logger, apiRouter)
	return router
}
