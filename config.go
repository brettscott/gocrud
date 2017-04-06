package main

import (
	"github.com/kelseyhightower/envconfig"
)

// AppConfig holds the application's configuration
type appConfig struct {
	Port          int    `envconfig:"PORT" required:"true"`
	Env           string `default:"local"`
	ComponentName string `envconfig:"COMPONENT_NAME" default:"gocrud"`

	MongoConnection string `envconfig:"MONGO_DB_CONNECTION" required:"true"`
	MongoSSLCert    string `envconfig:"MONGO_DB_SSL_CERT" default:""`
	MongoDBName     string `envconfig:"MONGO_DB_NAME" required:"true"`
}

func (c *appConfig) IsLocal() bool {
	return c.Env == "local"
}

func loadAppConfig() (*appConfig, error) {
	var config appConfig
	err := envconfig.Process("", &config)
	return &config, err
}
