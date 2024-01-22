package config

import (
	"fmt"
	"github.com/jackvonhouse/enrichment/pkg/log"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

type Database struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
	SSLMode      string
}

func (d Database) String() string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		d.Username, d.Password,
		d.Host, d.Port,
		d.DatabaseName,
		d.SSLMode,
	)
}

type ServerHTTP struct {
	Port int
}

type Config struct {
	Database Database
	Server   ServerHTTP
}

func New(
	configPath string,
	logger log.Logger,
) (Config, error) {

	configType := strings.TrimPrefix(filepath.Ext(configPath), ".")

	viper.SetConfigType(configType)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		logger.WithFields(map[string]any{
			"layer":       "config",
			"config_path": configPath,
		}).Warnf("error on reading config: %s", err)

		return Config{}, fmt.Errorf("error on reading config: %s", err)
	}

	postgresPrefix := "database.postgres"

	return Config{
		Database: Database{
			Username: viper.GetString(
				fmt.Sprintf("%s.username", postgresPrefix),
			),

			Password: viper.GetString(
				fmt.Sprintf("%s.password", postgresPrefix),
			),

			Host: viper.GetString(
				fmt.Sprintf("%s.host", postgresPrefix),
			),

			Port: viper.GetInt(
				fmt.Sprintf("%s.port", postgresPrefix),
			),

			DatabaseName: viper.GetString(
				fmt.Sprintf("%s.database_name", postgresPrefix),
			),

			SSLMode: viper.GetString(
				fmt.Sprintf("%s.ssl_mode", postgresPrefix),
			),
		},

		Server: ServerHTTP{
			Port: viper.GetInt("server.http.port"),
		},
	}, nil
}
