package main

import (
	log "github.com/Sirupsen/logrus"
	"time"
)

func Run(conf DictatorConfiguration, stop <-chan bool, finished chan<-bool) {
	var re Redis // Create a Redis Node
	err := re.Initialize(conf.Node.Name, conf.Node.Host, conf.Node.Port)
	if err != nil {
		log.Warn("Fail to initialize Redis node")
		finished <- true
	}

	var ze Elector // Create a ZK Elector
	err = ze.Initialize(conf.ZKHosts, conf.ServiceName, &re)
	if err != nil {
		log.Warn("Fail to initialize ZK Elector")
		finished <- true
	}

	// Run Elector
    go ze.Run()

	// Wait for the stop signal
	Loop:
	for {
		select {
		case hasToStop := <-stop:
			if hasToStop {
				log.Debug("Close Signal Received!")
			}else {
				log.Debug("Close Signal Received (but a strange false one)")
			}
			break Loop
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}

	ze.Destroy()

	finished <- true
}