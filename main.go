package main

import (
	"commerce/auth"
	"commerce/config"
	"commerce/models"
	"commerce/routes"
)

func main() {
	cfg := config.GenerateConfig()
	m := models.NewModels(
		models.WithMysql(cfg.Database.ConnectionString()),
		models.WithUserModel(),
		models.WithAccountsModel(),
	)
	jwt := auth.InitAuth(cfg.Secret)
	m.ApplyMigration()

	routes.Run(cfg.Port, m, cfg.Salt, jwt)
}
