// Package config to env
package config

import (
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

// Config struct to config env
type Config struct {
	Port     string `env:"PORT" envDefault:":8080" json:"port,omitempty"`
	Host     string `env:"HOST" envDefault:"localhost" json:"host,omitempty"`
	LogLevel string `env:"LOGLEVEL" envDefault:"debug" json:"loglevel,omitempty"`
	DBURL    string `env:"DBURL" envDefault:"postgres://andeisaldyun:e3cr3t@localhost:5432/catsDB" json:"dburl,omitempty"`

	HashSalt                    string `env:"HASHSALT" envDefault:"HAsh_salt" json:"hash_salt,omitempty"`
	AuthenticationKey           string `env:"AUTHENTICATIONKEY" envDefault:"authentication_key" json:"authentication_key,omitempty"`
	RefreshKey                  string `env:"REFRESHKEY" envDefault:"refresh_key" json:"refresh_key,omitempty"`
	AuthenticationTokenDuration int    `env:"TOKENDURATION" envDefault:"3600" json:"token_duration,omitempty"`
	RefreshTokenDuration        int    `env:"REFRESHTOKENDURATION" envDefault:"86400" json:"refresh-token-duration,omitempty"`
}

// New contract config
func New(conf *Config) {
	if err := env.Parse(conf); err != nil {
		log.Fatalf("Unable to parse config")
	}
}
