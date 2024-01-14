package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	DBConfig  DatabaseConfig
	SvcConfig ServiceConfig
	MailConf  MailConfig
}

type ServiceConfig struct {
	LogLevel          string `envconfig:"LOG_LEVEL" default:"info"`
	Port              uint   `envconfig:"PORT" default:"8080"`
	ShutdownWait      uint16 `envConfig:"SHUTDOWN_WAIT" default:"20"`
	HeaderReadTimeout uint16 `envConfig:"HEADER_READ_TIMEOUT" default:"20"`
	GinAccessLog      bool   `envconfig:"GIN_ACCESS_LOG" default:"false"`
}

type DatabaseConfig struct {
	DatabaseType string `envconfig:"DATABASE_TYPE" default:"mysql"`
	Username     string `envconfig:"USERNAME" default:"root"`
	Password     string `envconfig:"PASSWORD" default:"root"`
	Port         string `envconfig:"PORT" default:"3306"`
	DbName       string `envconfig:"DATABASE_NAME" default:"testing"`
	Url          string `envconfig:"URL" default:"127.0.0.1"`
	SslConfig    SSLConfig
}

type SSLConfig struct {
	Sslmode    string `envconfig:"SSLMODE" default:"disable"`
	MinTLS     string `envconfig:"MIN_TLS" default:"10"`
	RootCA     string `envconfig:"ROOT_CA" default:"test"`
	ServerCert string `envconfig:"SERVER_CERT" default:"test"`
	ClientCert string `envconfig:"CLIENT_CERT" default:"test"`
	ClientKey  string `envconfig:"CLIENT_KEY" default:"test"`
}

type MailConfig struct {
	Host     string `envconfig:"MAIL_HOST" default:"test"`
	Port     int    `envconfig:"MAIL_PORT" default:"test"`
	Username string `envconfig:"MAIL_USERNAME" default:"test"`
	Password string `envconfig:"MAIL_PASSWORD" default:"test"`
	Sender   string `envconfig:"MAIL_SENDER" default:"test"`
}

func NewConfig() (*Configuration, error) {
	var dbconfig DatabaseConfig
	if err := envconfig.Process("", &dbconfig); err != nil {
		return nil, fmt.Errorf("database configuration failed %v", err)
	}

	var svcConfig ServiceConfig
	if err := envconfig.Process("", &svcConfig); err != nil {
		return nil, fmt.Errorf("service configuration failed %v", err)
	}

	var mailConfig MailConfig
	if err := envconfig.Process("", &mailConfig); err != nil {
		return nil, fmt.Errorf("mail configuration failed %v", err)
	}

	return &Configuration{
		DBConfig:  dbconfig,
		SvcConfig: svcConfig,
		MailConf:  mailConfig,
	}, nil
}
