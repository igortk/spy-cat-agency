package main

import (
	"log"
	"spy-cat-agency/config"
	"spy-cat-agency/service"
	"spy-cat-agency/utils"

	"spy-cat-agency/database"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Failed load configuration: %v", err)
	}

	database.ConnectDB(&cfg.DataBaseConfig)

	utils.InitializeConfig(cfg.ExternalServicesConfig.TheCatAPIURL)

	service.Run(&cfg.HttpConfig)
}
