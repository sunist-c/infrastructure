package postgres

import (
	"errors"
	"fmt"
	"github.com/alioth-center/infrastructure/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const DriverNamePostgres = "postgres"

type postgresDatabase struct{}

func NewPostgresDatabase() adb.Database { return &postgresDatabase{} }

func (db *postgresDatabase) InitDriver(conf adb.Config) (driver *gorm.DB, err error) {
	if conf.DriverName == DriverNamePostgres {
		return gorm.Open(postgres.New(postgres.Config{
			DriverName: DriverNamePostgres,
			DSN:        db.dsn(conf.DataSource),
		}), db.config(conf.DriverOptions))
	}

	return nil, errors.New("unsupported postgres database driver")
}

func (db *postgresDatabase) dsn(dataSource map[string]any) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		dataSource["host"], dataSource["port"],
		dataSource["username"], dataSource["password"],
		dataSource["dbname"], dataSource["ssl_mode"], dataSource["location"],
	)
}

func (db *postgresDatabase) config(_ map[string]any) *gorm.Config {
	// todo: support postgres gorm config
	return &gorm.Config{}
}
