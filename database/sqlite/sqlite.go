package sqlite

import (
	"errors"
	"fmt"
	"github.com/alioth-center/infrastructure/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DriverNameSqlite = "sqlite"

type sqliteDatabase struct{}

func NewSqliteDatabase() adb.Database { return &sqliteDatabase{} }

func (db *sqliteDatabase) InitDriver(conf adb.Config) (driver *gorm.DB, err error) {
	if conf.DriverName == DriverNameSqlite {
		return gorm.Open(sqlite.Open(db.dsn(conf.DataSource)), db.config(conf.DriverOptions))
	}

	return nil, errors.New("unsupported sqlite database driver")
}

func (db *sqliteDatabase) dsn(dataSource map[string]any) string {
	return fmt.Sprintf("file:%s%s", dataSource["file"], dataSource["cache"])
}

func (db *sqliteDatabase) config(_ map[string]any) *gorm.Config {
	// todo: support sqlite gorm config
	return &gorm.Config{}
}
