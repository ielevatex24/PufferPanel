package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pufferpanel/pufferpanel/v2/config"
	"github.com/pufferpanel/pufferpanel/v2/models"
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
			ID: "1646495734",
			Migrate: func(db *gorm.DB) error {
				dialect := config.GetString("panel.database.dialect")
				// gorm can migrate other databases itself, only sqlite needs this migration done manually
				if dialect == "" || dialect == "sqlite3" {
					type NewSetting struct {
						Key   string `gorm:"type:varchar(100);primaryKey"`
						Value string `gorm:"type:text"`
					}
					if err := db.AutoMigrate(&NewSetting{}); err != nil {
						return err
					}
					var settings []config.Setting
					result := db.Find(&settings)
					if result.Error != nil {
						return result.Error
					}
					migrated := []NewSetting{}
					for _, v := range settings {
						migrated = append(migrated, NewSetting{
							Key: v.Key,
							Value: v.Value,
						})
					}
					result = db.Create(&migrated)
					if result.Error != nil {
						return result.Error
					}
					if err := db.Migrator().DropTable("settings"); err != nil {
						return err
					}
					if err := db.Migrator().RenameTable("new_settings", "settings"); err != nil {
						return err
					}
				}
				return nil
			},
			Rollback: nil,
		},
	})

	return m.Migrate()
}
