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

type S3Config struct {
	Region   string `json:"region"`
	Endpoint string `json:"endpoint"`
	Bucket   string `json:"bucket"`
}

type Service struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Config struct {
	Postgres     DBConfig `json:"postgres"`
	Tarantool    DBConfig `json:"tarantool"`
	Payment      string   `json:"payment"`
	Info         Service  `json:"info"`
	Auth         Service  `json:"auth"`
	Stream       Service  `json:"stream"`
	Subscription Service  `json:"subscription"`
	S3           S3Config `json:"aws_s3"`
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
	}
}
