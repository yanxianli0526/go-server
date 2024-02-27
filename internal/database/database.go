package database

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"meepShopTest/config"
	"meepShopTest/internal/model"
)

type GormDatabase struct {
	DB *gorm.DB
}

func initEngine(config *config.Config, logger *zap.Logger) *GormDatabase {
	var err error
	dsn := fmt.Sprintf("sslmode=%s host=%s port=%v dbname=%s password=%s user=%s", config.DB.Postgres.SSLMode, config.DB.Postgres.Host, config.DB.Postgres.Port, config.DB.Postgres.DBName, config.DB.Postgres.Password, config.DB.Postgres.User)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("connect database failed:", zap.String("addr", dsn), zap.Error(err))
	}

	return &GormDatabase{DB: db}
}

func migration(d GormDatabase) error {
	d.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	return d.DB.AutoMigrate(
		&model.User{},
	)
}

func GetPostgresCli(config *config.Config) *GormDatabase {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	db := initEngine(config, logger)
	err := migration(*db)
	if err != nil {
		logger.Fatal("postgres migration have err:", zap.Error(err))
	}

	return db
}
