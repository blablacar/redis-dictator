package main

import (
	log "github.com/Sirupsen/logrus"
	"strconv"
	"fmt"
	"net/http"
)


func HTTPDictatorStatus(w http.ResponseWriter, r *http.Request, ze *Elector) {
	if ze.Paused {
		fmt.Fprintf(w, "Dictator Status = Disable/Paused\n")
		log.Info("Dictator Status = Disable/Paused")
	}else{
		fmt.Fprintf(w, "Dictator Status = Enable\n")
		log.Info("Dictator Status = Enable")
	}	
}


func HTTPEnable(w http.ResponseWriter, r *http.Request, ze *Elector) {
	if ze.Paused {
		go ze.Run()
		ze.Paused = false
		log.Info("Unpause Dictator")
		fmt.Fprintf(w, "Dictator Status = Enable\n")
	}else{
		log.Warn("Dictator is already enabled")
	}
}

func HTTPDisable(w http.ResponseWriter, r *http.Request, ze *Elector) {
	if !ze.Paused {
		ze.Destroy()
		ze.Paused = true
		log.Info("Pause Dictator")
		fmt.Fprintf(w, "Dictator Status = Disable/Paused\n")
	}else{
		log.Warn("Dictator is already paused")
	}
}

func Run(conf DictatorConfiguration, stop <-chan bool, finished chan<-bool) {
	var re Redis // Create a Redis Node
	err := re.Initialize(conf.Node.Name, conf.Node.Host, conf.Node.Port, conf.Node.LoadingTimeout)
	if err != nil {
		log.WithError(err).Warn("Fail to initialize Redis node")
		finished <- true
	}

	var ze Elector // Create a ZK Elector
	err = ze.Initialize(conf.ZKHosts, conf.ServiceName, &re)
	if err != nil {
		log.WithError(err).Warn("Fail to initialize ZK Elector")
		finished <- true
	}

	// Run Elector
    go ze.Run()

    // http signals management
    http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
          HTTPDictatorStatus(w, r, &ze)
    })
    http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
          HTTPDictatorStatus(w, r, &ze)
    })
    http.HandleFunc("/enable", func(w http.ResponseWriter, r *http.Request) {
          HTTPEnable(w, r, &ze)
    })
    http.HandleFunc("/disable", func(w http.ResponseWriter, r *http.Request) {
          HTTPDisable(w, r, &ze)
    })
	go http.ListenAndServe(":" + strconv.Itoa(conf.HttpPort), nil)

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
		}
	}

	ze.Destroy()

	finished <- true
}