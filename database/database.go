package database

import (
	"fmt"
	"time"

	"github.com/horlakz/go-auth/constants"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DatabaseFacade *gorm.DB

type DatabaseInterface interface {
	Connection() *gorm.DB
}

type connection struct {
	database *gorm.DB
}

func StartDatabaseClient(env constants.Env) DatabaseInterface {
	dsn := "host=" + env.DB_HOST + " user=" + env.DB_USER + " password=" + env.DB_PASSWORD + " dbname=" + env.DB_NAME + " port=" + env.DB_PORT

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true})

	if err != nil {
		fmt.Print(err)
	}

	sqlDb, err := database.DB()

	sqlDb.SetMaxIdleConns(5)
	sqlDb.SetMaxOpenConns(10)
	sqlDb.SetConnMaxLifetime(time.Hour)

	if err != nil {
		fmt.Print(err)
	}

	fmt.Println("Database connection is successful")

	DatabaseFacade = database

	return &connection{
		database: database,
	}
}

func (conn connection) Connection() *gorm.DB {
	return conn.database
}
