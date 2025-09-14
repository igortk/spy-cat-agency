package database

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"spy-cat-agency/config"
	"spy-cat-agency/models"
)

var Session *gorm.DB

func ConnectDB(dbCfg *config.DataBaseConfig) {
	var err error
	
	Session, err = gorm.Open(postgres.Open(dbCfg.Host), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Info("Database connection established!")

	if err := Session.AutoMigrate(&models.Cat{}, &models.Mission{}, &models.Target{}); err != nil {
		log.Errorf("Failed auto migrate: %v", err)
	}
}
