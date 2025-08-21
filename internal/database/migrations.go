package database

import (
	"fmt"

	"github.com/andro-kes/SubAggr/internal/models"
	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Migrate() error {
	m := gormigrate.New(DB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "20250101_init_subs",
			Migrate: func(tx *gorm.DB) error {
				type Subs struct {
					models.Subs
				}
				if err := tx.AutoMigrate(&Subs{}); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("subs")
			},
		},
		{
			ID: "20250101_indexes",
			Migrate: func(tx *gorm.DB) error {
				// Ensure composite unique index and regular indexes exist
				if err := tx.Migrator().CreateIndex(&models.Subs{}, "uniq_user_service"); err != nil {
					return err
				}
				if err := tx.Migrator().CreateIndex(&models.Subs{}, "user_id"); err != nil {
					return err
				}
				if err := tx.Migrator().CreateIndex(&models.Subs{}, "service_name"); err != nil {
					return err
				}
				if err := tx.Migrator().CreateIndex(&models.Subs{}, "start_date"); err != nil {
					return err
				}
				if err := tx.Migrator().CreateIndex(&models.Subs{}, "end_date"); err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				// GORM migrator drops indexes by name
				_ = tx.Migrator().DropIndex(&models.Subs{}, "uniq_user_service")
				_ = tx.Migrator().DropIndex(&models.Subs{}, "user_id")
				_ = tx.Migrator().DropIndex(&models.Subs{}, "service_name")
				_ = tx.Migrator().DropIndex(&models.Subs{}, "start_date")
				_ = tx.Migrator().DropIndex(&models.Subs{}, "end_date")
				return nil
			},
		},
	})

	if err := m.Migrate(); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	return nil
}

