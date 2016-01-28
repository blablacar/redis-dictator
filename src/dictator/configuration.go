package main

import (
	log "github.com/Sirupsen/logrus"
	"encoding/json"
	"io/ioutil"
)

type NodeConfiguration struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int `json:"port"`
}

type DictatorConfiguration struct {
	ServiceName string `json:"svc_name"`
	LogLevel    string `json:"log_level"`
	ZKHosts     []string `json:"zk_hosts"`
	Node        NodeConfiguration `json:"node"`
}

func NewDictatorConfiguration() DictatorConfiguration {
	log.Debug("Initialize configuration")

	return DictatorConfiguration{
		LogLevel: "INFO",
		ZKHosts: []string{"localhost:2181"},
		Node: NodeConfiguration{
			Name: "local",
			Host:"localhost",
			Port: 6379,
		},
	}
}

func (d *DictatorConfiguration) ReadConfigurationFile(configFilePath string) error {
	if configFilePath == "" {
		return nil
	}
	log.WithField("file", configFilePath).Debug("Reading configuration file")

	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, d)
	if err != nil {
		return err
	}


	return nil
}