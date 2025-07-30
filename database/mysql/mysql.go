package mysql

import (
	"errors"
	"fmt"
	"github.com/alioth-center/infrastructure/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DriverNameMysql   = "mysql"
	DriverNameMariadb = "mariadb"
)

type mysqlDatabase struct{}

func NewMySqlDatabase() adb.Database { return &mysqlDatabase{} }

func (db *mysqlDatabase) InitDriver(conf adb.Config) (driver *gorm.DB, err error) {
	switch conf.DriverName {
	case DriverNameMysql:
		return gorm.Open(mysql.New(mysql.Config{
			DriverName:        DriverNameMysql,
			DSN:               db.dsn(conf.DataSource),
			DefaultStringSize: conf.StringSize,
		}), db.config(conf.DriverOptions))
	case DriverNameMariadb:
		return gorm.Open(mysql.New(mysql.Config{
			DriverName:                DriverNameMysql,
			DSN:                       db.dsn(conf.DataSource),
			SkipInitializeWithVersion: false,
			DefaultStringSize:         conf.StringSize,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
		}), db.config(conf.DriverOptions))
	default:
		return nil, errors.New("unsupported mysql database driver")
	}
}

func (db *mysqlDatabase) dsn(dataSource map[string]any) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s",
		dataSource["username"], dataSource["password"],
		dataSource["host"], dataSource["port"],
		dataSource["database"], dataSource["charset"], dataSource["location"],
	)
}

func (db *mysqlDatabase) config(_ map[string]any) *gorm.Config {
	// todo: support mysql gorm config
	return &gorm.Config{}
}
