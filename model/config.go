package model

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ServiceSetting ServiceSetting
	SQLSetting     SQLSetting
}

func ConfigFromFile(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return ConfigFromJSON(data)
}

func ConfigFromJSON(data []byte) (*Config, error) {
	conf := &Config{}
	err := json.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

type ServiceSetting struct {
	APIEndpoint string
}

type SQLSetting struct {
	DataSource         string
	MaxIdleConnections int
	MaxOpenConnections int
}
