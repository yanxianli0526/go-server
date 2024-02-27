package server

import (
	"context"
	"io"
	"log"
	"meepShopTest/config"
	"meepShopTest/internal/database"
	"meepShopTest/internal/router"
	"net/http"
	"os"
	"os/signal"

	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Run(config *config.Config, db *database.GormDatabase) {
	if config.Gin.Mode == gin.DebugMode || config.Gin.Mode == gin.TestMode || config.Gin.Mode == gin.ReleaseMode {
		gin.SetMode(config.Gin.Mode)
	}
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ginEngine := gin.Default()
	ginEngine = router.New(ginEngine, logger, db)
	ginEngine.NoRoute(notFound)

	srv := &http.Server{
		Addr:    ":" + config.Server.Port,
		Handler: ginEngine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	gracefulShutdown(srv, config.Server.ShutdownTimeoutSec)
}

func gracefulShutdown(srv *http.Server, shutdownTimeoutSec int) {
	quit := make(chan os.Signal, 3)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeoutSec)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}

func notFound(c *gin.Context) {
	if c.Request != nil && c.Request.Method == http.MethodPost {
		jsonData, err := io.ReadAll(c.Request.Body)
		if err == nil {
			log.Printf("jsonData: %s\n", jsonData)
		}
	}

	c.AbortWithStatus(http.StatusNotFound)
}
