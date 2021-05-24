package config

import (
	"Redioteka/internal/pkg/utils/log"
	"encoding/json"
	"io/ioutil"
)

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBName   string `json:"db,omitempty"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Config struct {
	Postgres  *DBConfig `json:"postgres"`
	Tarantool *DBConfig `json:"tarantool"`
	Payment   string    `json:"payment"`
}

var config = &Config{}

func Get() *Config {
	return config
}

func init() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Log.Error(err)
		return
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		log.Log.Error(err)
		return
	}
}
