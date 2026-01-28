package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/shadowbane/go-logger"
	"github.com/spf13/viper"
)

type Config struct {
	port string

	dbUser string
	dbPass string
	dbPort string
	dbHost string
	dbName string
}

func Get() *Config {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	// Initialize logger
	logger.Init(logger.LoadEnvForLogger())

	return &Config{
		port:   viper.GetString("PORT"),
		dbUser: viper.GetString("DB_USERNAME"),
		dbPass: viper.GetString("DB_PASSWORD"),
		dbPort: viper.GetString("DB_PORT"),
		dbHost: viper.GetString("DB_HOST"),
		dbName: viper.GetString("DB_DATABASE"),
	}
}

func (c *Config) GetPort() string {
	return c.port
}

func (c *Config) GetAddr() string {
	return "0.0.0.0:" + c.port
}

func (c *Config) GetPSQLConnectionString() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		url.QueryEscape(c.dbUser),
		url.QueryEscape(c.dbPass),
		c.dbHost,
		c.dbPort,
		c.dbName,
	)
}
