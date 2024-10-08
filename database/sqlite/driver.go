package sqlite

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alioth-center/infrastructure/database"
	"github.com/alioth-center/infrastructure/exit"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const DriverName = "sqlite"

type sqliteDb struct {
	database.BaseDatabaseImplement
	database.BaseDatabaseImplementV2

	initialized bool
}

func (s *sqliteDb) Init(options database.Options) error {
	// 防止重复初始化
	if s.initialized {
		return nil
	}
	s.initialized = true

	// 初始化日志器
	if options.Logger == nil {
		options.Logger = logger.Default()
	}
	dataSource := options.DataSource
	s.Logger = options.Logger
	options.Logger.Info(logger.NewFields().WithMessage("start open sqliteDb database").WithData(dataSource))

	// 连接数据库
	db, openErr := gorm.Open(sqlite.Open(dataSource))
	if openErr != nil {
		if !os.IsNotExist(openErr) {
			err := fmt.Errorf("open sqliteDb database error: %w", openErr)
			s.Logger.Error(logger.NewFields().WithMessage(err.Error()).WithData(dataSource))
			return err
		} else if _, createErr := os.Create(filepath.Join("data", dataSource)); createErr != nil {
			err := fmt.Errorf("create sqliteDb database error: %w", createErr)
			s.Logger.Error(logger.NewFields().WithMessage(err.Error()).WithData(dataSource))
			return err
		}
	}
	db.Logger = database.NewDBLogger(options.Logger)

	// 设置数据库连接池
	sqlDb, dbe := db.DB()
	if dbe != nil {
		return fmt.Errorf("get sqliteDb database error: %w", dbe)
	}
	s.BaseDatabaseImplement.ParseDatabaseOptions(sqlDb, options)
	s.BaseDatabaseImplement.SetRandCommand("random()")
	s.BaseDatabaseImplement.SetDriverName(DriverName)

	// 连接成功
	s.BaseDatabaseImplement.Db, s.BaseDatabaseImplementV2.Db = db, db
	s.Logger.Info(logger.NewFields().WithMessage("successfully open sqliteDb database").WithData(dataSource))

	// 注册退出事件
	exit.RegisterExitEvent(func(_ os.Signal) {
		_ = sqlDb.Close()
		fmt.Println("closed sqlite database")
	}, "CLOSE_SQLITE_DB_CONN")
	return nil
}

// NewSqliteDb creates a new sqlite database instance.
//
// Deprecated: Use NewSQLiteV2 instead.
func NewSqliteDb(config Config, models ...any) (db database.Database, err error) {
	rdb := &sqliteDb{}
	if initErr := rdb.Init(convertConfigToOptions(config)); initErr != nil {
		return nil, fmt.Errorf("init sqliteDb database error: %w", initErr)
	} else if migrateErr := rdb.Migrate(models...); migrateErr != nil {
		return nil, fmt.Errorf("migrate sqliteDb database error: %w", migrateErr)
	} else {
		return rdb, nil
	}
}

// NewSQLiteV2 creates a new sqlite database instance.
func NewSQLiteV2(config Config, models ...any) (db database.DatabaseV2, err error) {
	rdb := &sqliteDb{}
	if initErr := rdb.Init(convertConfigToOptions(config)); initErr != nil {
		return nil, fmt.Errorf("init sqliteDb database error: %w", initErr)
	}

	if migrateErr := rdb.Migrate(models...); migrateErr != nil {
		return nil, fmt.Errorf("migrate sqliteDb database error: %w", migrateErr)
	}

	return rdb, nil
}

func NewWithLogger(config Config, logger logger.Logger, models ...any) (db database.DatabaseV2, err error) {
	rdb := &sqliteDb{}
	opts := convertConfigToOptions(config)
	opts.Logger = logger
	if initErr := rdb.Init(opts); initErr != nil {
		return nil, fmt.Errorf("init sqliteDb database error: %w", initErr)
	}

	if migrateErr := rdb.Migrate(models...); migrateErr != nil {
		return nil, fmt.Errorf("migrate sqliteDb database error: %w", migrateErr)
	}

	return rdb, nil
}
