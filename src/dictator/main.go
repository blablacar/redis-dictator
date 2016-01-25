package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"flag"
	"time"
	"os/signal"
)

// Manage OS Signal, only for shutdown purpose
// When termination signal is received, we send a message to a chan
func manageSignal(c <-chan os.Signal, stop chan<-bool) {
    for {
		select {
		case _signal := <-c:
			if _signal == os.Kill {
				log.Debug("Dictator: Kill Signal Received")
			}
			if _signal == os.Interrupt {
				log.Debug("Dictator: Interrupt Signal Received")
			}
			log.Debug("Dictator: Shutdown Signal Received")
			stop <- true
			break
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

// All the command line arguments are managed inside this function
func initFlags() (string) {
	// The configuration filename
	var configurationFileName = flag.String("config", "./dictator.json.conf", "the complete filename of the configuration file")
	// Parse all command line options
	flag.Parse()
	return *configurationFileName
}

func initConfiguration(configurationFileName string) (DictatorConfiguration) {
	log.Debug("Starting config file parsing")
	dictatorConfiguration , err := OpenConfiguration(configurationFileName)
	if err != nil {
		// If an error is raised when parsing configuration file
		// the configuration object can be either empty, either incomplete
		log.Fatal("Unable to load Configuration (",err,")")
		// So the configuration is incomplete, exit the program now
		os.Exit(1)
	}
	log.Debug("Config file parsed")
	return dictatorConfiguration
}

func main() {
	log.Info("Dictator: Starting")
    
    // Read Options
	configurationFileName := initFlags()

	// Load the configuration from config file. If something wrong, the full process is stopped inside the function
	dictatorConfiguration := initConfiguration(configurationFileName)

	c := make(chan os.Signal,1)
	signal.Notify(c, os.Interrupt, os.Kill)
	log.Debug("Dictator: Signal Channel notification setup done")

	stop := make(chan bool)
	go manageSignal(c,stop)
	log.Debug("Dictator: Signal Management Started")

    finished:= make(chan bool)
	go Run(dictatorConfiguration, stop,finished)
	log.Debug("Dictator: Go routine launched")

	log.Debug("Dictator: Waiting for main process to Stop")
	isFinished := <-finished
	if (isFinished) {
		log.Debug("Nerve: Main routine closed correctly")
	}else {
		log.Warn("Nerve: Main routine closed incorrectly")
	}
	log.Info("Dictator: Shutdown")
}