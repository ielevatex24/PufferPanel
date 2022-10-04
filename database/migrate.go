package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"gorm.io/gorm"
)

func migrate(dbConn *gorm.DB) error {

	m := gormigrate.New(dbConn, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "1626910428",
			Migrate: func(db *gorm.DB) error {
				_ = db.Migrator().DropIndex(&models.Server{}, "uix_servers_name")
				return nil
			},
			Rollback: nil,
		},
		{
			ID: "1658926619",
			Migrate: func(db *gorm.DB) error {
				return db.Create(&models.TemplateRepo{
					Name:   "community",
					Url:    "https://github.com/pufferpanel/templates",
					Branch: "v2",
				}).Error
			},
		},
	})

	return m.Migrate()
}
