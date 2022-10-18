package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"gorm.io/gorm"
)

func migrate(dbConn *gorm.DB) error {
	//first step is for nodes, we need to drop columns, but we to handle that first before we can migrate the models
	m := gormigrate.New(dbConn, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "1665758386",
			Migrate: func(db *gorm.DB) (err error) {
				migrator := db.Migrator()

				model := &models.Node{}
				for _, v := range []string{"private_host", "public_port", "private_port"} {
					if migrator.HasColumn(model, v) {
						err = migrator.DropColumn(model, v)
						if err != nil {
							return
						}
					}
				}
				return
			},
			Rollback: nil,
		},
	})

	err := m.Migrate()
	if err != nil {
		return err
	}

	dbObjects := []interface{}{
		&models.Node{},
		&models.Server{},
		&models.User{},
		&models.Permissions{},
		&models.Client{},
		&models.UserSetting{},
		&models.Session{},
		&models.TemplateRepo{},
	}

	for _, v := range dbObjects {
		if err = dbConn.AutoMigrate(v); err != nil {
			return err
		}
	}

	m = gormigrate.New(dbConn, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "1626910428",
			Migrate: func(db *gorm.DB) error {
				migrator := db.Migrator()

				if migrator.HasIndex(&models.Server{}, "uix_servers_name") {
					return migrator.DropIndex(&models.Server{}, "uix_servers_name")
				}
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
		{
			ID: "1665609381",
			Migrate: func(db *gorm.DB) error {
				var nodes []*models.Node
				err := db.Find(&nodes).Error
				if err != nil {
					return err
				}

				var local *models.Node
				for _, v := range nodes {
					if v.IsLocal() {
						local = v
					}
				}

				if local == nil {
					return nil
				}

				err = db.Table("servers").Where("node_id = ?", local.ID).Update("node_id", 0).Error
				if err != nil {
					return err
				}
				err = db.Delete(local).Error
				if err != nil {
					return err
				}

				return nil
			},
		},
	})

	return m.Migrate()
}
