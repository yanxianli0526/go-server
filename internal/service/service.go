package service

import (
	"context"

	"meepShopTest/internal/database"
)

type Service struct {
	ctx context.Context
	db  *database.GormDatabase
}

func New(ctx context.Context, db *database.GormDatabase) Service {
	service := Service{
		ctx: ctx,
		db:  db,
	}
	return service
}
