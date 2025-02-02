package migrations

import (
	"fmt"

	"github.com/TIM-DEBUG-ProjectSprintBatch3/go-fiber-template/src/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate() {
	if config.GetAutoMigrate() {
		migrated, err := migrate.New(config.GetLocationMigrate(), config.GetDBConnectionMigrate())
		if err != nil {
			message := fmt.Sprintf("Error creating migrate instance: %v", err)
			panic(message)
		}

		if err := migrated.Up(); err != nil {
			if err.Error() == "no change" {
				fmt.Println("No new migrations to apply.")
			} else {
				if err != nil {
					message := fmt.Sprintf("Error creating migrate instance: %v", err)
					panic(message)
				}
			}
		}
		fmt.Println("Migration completed successfully!")
	}
}
