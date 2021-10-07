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
}

// New contract config
func New(conf *Config) {
	if err := env.Parse(conf); err != nil {
		log.Fatalf("Unable to parse config")
	}
}
