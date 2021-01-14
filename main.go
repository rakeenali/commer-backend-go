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
		models.WithUserRoleModel(),
		models.WithTagsModel(),
		models.WithItemsModel(),
	)
	jwt := auth.InitAuth(cfg.Secret)
	m.ApplyMigration()

	r := routes.NewRouter(
		routes.WithModel(m),
		routes.WithMiddlewares(jwt),
		routes.WithUserRouter(jwt, cfg.Salt),
	)
	r.Run(cfg.Port)
}
