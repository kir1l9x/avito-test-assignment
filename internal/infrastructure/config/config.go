package config

import (
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

type DBConfig struct {
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT,required"`

	DBMaxCons     int           `env:"DB_MAX_CONS" envDefault:"10"`
	DBMinCons     int           `env:"DB_MIN_CONS" envDefault:"2"`
	DBConLifetime time.Duration `env:"DB_CON_LIFETIME" envDefault:"1h"`
	DBConnTimeout time.Duration `env:"DB_CONN_TIMEOUT" envDefault:"10s"`
}

func LoadDBConfig() (*DBConfig, error) {
	dbCfg := &DBConfig{}
	if err := env.Parse(dbCfg); err != nil {
		return nil, err
	}

	return dbCfg, nil
}

func (dbCfg *DBConfig) GetDBDSNPostgres() string {
	builder := strings.Builder{}
	builder.WriteString("postgresql://")
	builder.WriteString(dbCfg.DBUser)
	builder.WriteString(":")
	builder.WriteString(dbCfg.DBPassword)
	builder.WriteString("@")
	builder.WriteString(dbCfg.DBHost)
	builder.WriteString(":")
	builder.WriteString(dbCfg.DBPort)
	builder.WriteString("/")
	builder.WriteString(dbCfg.DBName)
	builder.WriteString("?sslmode=disable")

	return builder.String()
}

type AppConfig struct {
	HTTPPort string `env:"HTTP_PORT" envDefault:":8080"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

func LoadAppConfig() (*AppConfig, error) {
	appCfg := &AppConfig{}
	if err := env.Parse(appCfg); err != nil {
		return nil, err
	}

	if !strings.HasPrefix(appCfg.HTTPPort, ":") {
		appCfg.HTTPPort = ":" + appCfg.HTTPPort
	}

	return appCfg, nil
}
