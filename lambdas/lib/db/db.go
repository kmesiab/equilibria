package db

import (
	"errors"
	"fmt"

	mysql2 "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/log"
	"github.com/kmesiab/equilibria/lambdas/lib/utils"
)

var globalDB *gorm.DB

// Init initializes the database connection using the provided configuration.
func Init(config *config.Config) (*gorm.DB, error) {

	// Format the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseHost,
		config.DatabaseName,
	)

	// Open the database with the MySQL driver
	return gorm.Open(mysql.Open(dsn), &gorm.Config{

		Logger: logger.Default.LogMode(logger.LogLevel(config.LogLevel)),
	})
}

// Get returns the database connection instance.
func Get(config *config.Config) *gorm.DB {
	var err error

	if globalDB == nil {

		globalDB, err = Init(config)

		// If we don't have a database connection, panic
		if err != nil {
			panic(fmt.Sprintf("failed to initialize database: %s\n", err))
		}

		// If we can't reach the database, panic
		if utils.PingDatabase(globalDB) != nil {
			msg := fmt.Sprintf("failed to ping database: %s\n", err)
			panic(msg)
		} else {
			log.New("Database connected").Log()
		}
	}

	return globalDB
}

func IsDuplicateEntryError(err error) bool {
	var mysqlErr *mysql2.MySQLError

	if ok := errors.As(err, &mysqlErr); ok {
		return mysqlErr.Number == 1062
	}

	return false
}
