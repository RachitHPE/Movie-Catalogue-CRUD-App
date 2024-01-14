package db

import (
	"catalogue-app/internal/config"
	"catalogue-app/internal/pkg/log"
	"context"
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBClient struct {
	dbConfig *config.Configuration
}

func (client DBClient) DBInit() (*gorm.DB, error) {
	address := client.dbConfig.DBConfig.Url

	if client.dbConfig.DBConfig.Port != "" {
		address += ":" + client.dbConfig.DBConfig.Port
	}

	dbInfo := client.dbConfig.DBConfig.Username + ":" +
		client.dbConfig.DBConfig.Password +
		"@tcp(" + address + ")/" +
		client.dbConfig.DBConfig.DbName +
		"?charset=utf8mb4&parseTime=True&loc=Local"

	if client.dbConfig.DBConfig.SslConfig.Sslmode != "disable" {
		// use host machine's root CAs to verify
		if client.dbConfig.DBConfig.SslConfig.Sslmode == "require" {
			dbInfo += "&tls=true"
		}

		// perform comprehensive SSL/TLS certificate validation using
		// certificate signed by a recognized CA or by a self-signed certificate
		if client.dbConfig.DBConfig.SslConfig.Sslmode == "verify-ca" || client.dbConfig.DBConfig.SslConfig.Sslmode == "verify-full" {
			dbInfo += "&tls=custom"
			err := client.InitTLSMySQL()
			if err != nil {
				log.Errorf(context.Background(), "failed initializing db")

				return nil, err
			}
		}
	}

	sqlDB, err := sql.Open("mysql", dbInfo)
	if err != nil {
		log.Errorf(context.Background(), "failed opening sql")

		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)   // max number of connections in the idle connection pool
	sqlDB.SetMaxOpenConns(2)    // max number of open connections in the database
	sqlDB.SetConnMaxLifetime(1) // max amount of time a connection may be reused

	dbVal, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(0)),
	})
	if err != nil {
		log.Errorf(context.Background(), "failed opening db")

		return nil, err
	}
	// Only for debugging
	if err == nil {
		fmt.Println("DB connection successful!")
	}

	return dbVal, nil
}
