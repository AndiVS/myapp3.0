// Package config to env
package config

import (
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

// Config struct to config env
type Config struct {
	Port          string `env:"PORT" envDefault:":8081" json:"port,omitempty"`
	Host          string `env:"HOST" envDefault:"postgres" json:"host,omitempty"`
	LogLevel      string `env:"LOGLEVEL" envDefault:"debug" json:"loglevel,omitempty"`
	DBURLPOSTGRES string `env:"DBURLPOSTGRES" envDefault:"postgres://andeisaldyun:e3cr3t@postgres:5432/catsDB" json:"dburlpostgres,omitempty"`
	DBURLMONGO    string `env:"DBURLMONGO" envDefault:"mongodb://andeisaldyun:e3cr3t@172.18.0.2:27017" json:"dburlmongo,omitempty"`
	DBURL         string `env:"DBURL" envDefault:"" json:"dburl,omitempty"`

	Server string `env:"SERVER" envDefault:"echo" json:"server"`
	Broker string `env:"BROKER" envDefault:"kafka" json:"broker"`

	System     string `env:"SYSTEM" envDefault:"postgres"`
	DBUser     string `env:"DB_USER" envDefault:"andeisaldyun"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"e3cr3t"`
	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     int    `env:"DB_PORT" envDefault:"5432"`
	DBName     string `env:"DB_NAME" envDefault:"catsDB"`

	/*System     string `env:"SYSTEM" envDefault:"mongodb"`
	DBUser     string `env:"DB_USER" envDefault:"andeisaldyun"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"e3cr3t"`
	DBHost     string `env:"DB_HOST" envDefault:"172.18.0.2"`
	DBPort     int    `env:"DB_PORT" envDefault:"27017"`
	DBName     string `env:"DB_NAME" envDefault:"catsDB"`*/

	HashSalt                    string `env:"HASHSALT" envDefault:"HAsh_salt" json:"hash_salt,omitempty"`
	AuthenticationKey           string `env:"AUTHENTICATIONKEY" envDefault:"authentication_key" json:"authentication_key,omitempty"`
	RefreshKey                  string `env:"REFRESHKEY" envDefault:"refresh_key" json:"refresh_key,omitempty"`
	AuthenticationTokenDuration int    `env:"TOKENDURATION" envDefault:"36000" json:"token_duration,omitempty"`
	RefreshTokenDuration        int    `env:"REFRESHTOKENDURATION" envDefault:"864000" json:"refresh-token-duration,omitempty"`
}

// New contract config
func New(conf *Config) {
	if err := env.Parse(conf); err != nil {
		log.Fatalf("Unable to parse config")
	}
}
