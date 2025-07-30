package sqlserver

import (
	"errors"
	"fmt"
	"github.com/alioth-center/infrastructure/database"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

const DriverNameSqlServer = "sqlserver"

type sqlserverDatabase struct{}

func NewSqlserverDatabase() adb.Database { return &sqlserverDatabase{} }

func (db *sqlserverDatabase) InitDriver(conf adb.Config) (driver *gorm.DB, err error) {
	if conf.DriverName == DriverNameSqlServer {
		return gorm.Open(sqlserver.Open(db.dsn(conf.DataSource)), db.config(conf.DriverOptions))
	}

	return nil, errors.New("unsupported sqlserver database driver")
}

func (db *sqlserverDatabase) dsn(dataSource map[string]any) string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		dataSource["username"], dataSource["password"],
		dataSource["host"], dataSource["port"], dataSource["database"],
	)
}

func (db *sqlserverDatabase) config(_ map[string]any) *gorm.Config {
	// todo: support sqlserver gorm config
	return &gorm.Config{}
}
