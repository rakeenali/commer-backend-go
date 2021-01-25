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
		models.WithUserBalanceModels(),
		models.WithOrders(),
	)
	jwt := auth.InitAuth(cfg.Secret)
	// m.DestroyTables()
	m.ApplyMigration()

	r := routes.NewRouter(
		routes.WithModel(m),
		routes.WithMiddlewares(jwt),
		routes.WithUserRouter(jwt, cfg.Salt),
		routes.WithTags(),
		routes.WithItemsRouter(),
		routes.WithOrdersRouter(),
	)
	r.Run(cfg.Port)
}
