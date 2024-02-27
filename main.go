/*
Copyright © 2024 yanxian <billyli@22451@gmail.com>
*/
package main

import (
	"meepShopTest/config"
	"meepShopTest/internal/database"
	"meepShopTest/internal/server"
)

func main() {
	cfg := config.New()
	postgresClient := database.GetPostgresCli(cfg)
	server.Run(cfg, postgresClient)
}
