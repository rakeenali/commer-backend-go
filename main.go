package main

import (
	"commerce/config"
	"commerce/models"
	"commerce/routes"
)

func main() {
	cfg := config.GenerateConfig()
	m := models.NewModels(
		models.WithMysql(cfg.Database.ConnectionString()),
		models.WithUserModel(),
	)
	// m.ApplyMigration()

	routes.Run(cfg.Port, m, cfg.Salt)
}
