package main

import (
	"encoding/json"
	"io/ioutil"
	"k8s.io/kubernetes/Godeps/_workspace/src/github.com/Sirupsen/logrus"
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
	logrus.Debug("Initialize configuration")

	return DictatorConfiguration{
		ServiceName:"local",
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
	logrus.WithField("file", configFilePath).Debug("Reading configuration file")

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