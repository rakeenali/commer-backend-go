package main

import (
	"commerce/config"
	"commerce/models"
	"commerce/routes"
	"fmt"
)

func main() {
	cfg := config.GenerateConfig()
	m := models.NewModels(
		models.WithMysql(cfg.Database.ConnectionString()),
		models.WithUserModel(),
	)
	fmt.Println(m)

	routes.Run(cfg.Port, m)
}
