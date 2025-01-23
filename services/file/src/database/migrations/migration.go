package migrations

import (
	"github.com/TimDebug/TutupLapak/File/src/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate() error {
	if config.GetAutoMigrate() {
		migrated, err := migrate.New(config.GetLocationMigrate(), config.GetDBConnectionMigrate())
		if err != nil {
			return err
		}

		if err := migrated.Up(); err != nil {
			if err.Error() == "no change" {
			} else {
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
