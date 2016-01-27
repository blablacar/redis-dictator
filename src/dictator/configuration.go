package main

import (
	"encoding/json"
	"io/ioutil"
)

type ElectorConfiguration struct {
	ZKHosts []string `json:"zk_hosts"`
	CheckInterval int `json:"check_interval"`
}

type NodeConfiguration struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int `json:"port"`
}

type DictatorConfiguration struct {
	ServiceName string `json:"svc_name"`
	LogLevel string `json:"log_level"`
	Elector ElectorConfiguration `json:"elector"`
	Node NodeConfiguration `json:"node"`
}

// Open Dictator configuration file, and parse it's JSON content
// return a full configuration object and an error
// if the error is different of nil, then the configuration object is empty
// if error is equal to nil, all the JSON content of the configuration file is loaded into the object
func OpenConfiguration(fileName string) (DictatorConfiguration, error) {
	var dictatorConfiguration DictatorConfiguration

	// Open and read the configuration file
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		// If there is an error with opening or reading the configuration file, return the error, and an empty configuration object
		return dictatorConfiguration, err
	}

	// Trying to convert the content of the configuration file (theoriticaly in JSON) into a configuration object
	err = json.Unmarshal(file, &dictatorConfiguration)
	if err != nil {
		// If there is an error in decoding the JSON entry into configuration object, return a partialy unmarshalled object, and the error
		return dictatorConfiguration, err
	}

	return dictatorConfiguration, nil
}
