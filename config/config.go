package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// MysqlConfig configuration for database
type MysqlConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// ConnectionString will create default mysql connection
func (m MysqlConfig) ConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", m.User, m.Password, m.Host, m.Name)
}

// DefaultConfig for api
type DefaultConfig struct {
	Port     int         `json:"port"`
	Env      string      `json:"env"`
	Salt     string      `json:"salt"`
	Database MysqlConfig `json:"database"`
}

// GenerateConfig will generate config for api
func GenerateConfig() *DefaultConfig {
	var config DefaultConfig
	data, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return &config
}
