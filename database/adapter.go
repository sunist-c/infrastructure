package adb

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type Database interface {
	InitDriver(conf Config) (driver *gorm.DB, err error)
}

type Config struct {
	StringSize         uint
	MaxIdleConnections uint
	MaxOpenConnections uint
	ConnectionLifeTime uint
	DriverName         string
	DriverOptions      map[string]any
	DataSource         map[string]any
}

func ConnectDatabase(db Database, conf Config, logger logger.Interface) (driver *gorm.DB, err error) {
	if db == nil {
		return nil, errors.New("database is nil")
	}

	// initialize gorm driver and gorm logger
	driver, err = db.InitDriver(conf)
	if err != nil {
		return nil, err
	}
	driver.Logger = logger

	// initialize database connection pool options
	rawDB, getErr := driver.DB()
	if getErr != nil {
		return nil, getErr
	}
	rawDB.SetMaxIdleConns(int(conf.MaxIdleConnections))
	rawDB.SetMaxOpenConns(int(conf.MaxOpenConnections))
	rawDB.SetConnMaxLifetime(time.Duration(conf.ConnectionLifeTime) * time.Second)

	return driver, nil
}
