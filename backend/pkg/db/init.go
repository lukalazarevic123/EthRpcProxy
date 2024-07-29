package db

import (
	"backend/config"
	"backend/pkg/db/model"
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var db *gorm.DB

func Init(cfg *config.DBConfig) (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Dbname, cfg.SslMode)
	sqlDB, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)
	newLogger.Info(context.Background(), "aaa")
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB.ConnPool,
	}), &gorm.Config{})
	fmt.Println("Created ", gormDB)

	tables := []interface{}{&model.HolderEntity{}}

	gormDB.AutoMigrate(tables...)
	db = gormDB

	return gormDB, err
}
