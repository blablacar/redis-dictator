package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"flag"
	"time"
	"os/signal"
	"fmt"
)

var (
	Version = "No Version Defined"
	BuildTime = "1970-01-01_00:00:00_UTC"
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
func initFlags() (string, string) {
	// The Log Level, from the Sirupsen/logrus level
	var logLevel = flag.String("log-level", "INFO", "A value to choose between [DEBUG INFO WARN FATAL], can be overriden by config file")
	// The configuration filename
	var configurationFileName = flag.String("config", "./dictator.json.conf", "the complete filename of the configuration file")
	// The version flag
	var version = flag.Bool("version", false, "Display version and exit")
	// Parse all command line options
	flag.Parse()
	if *version {
		printVersion()
		os.Exit(0)
	}
	return *logLevel, *configurationFileName
}

// Set the Logrus global log level
// Converted from a configuration string
func setLogLevel(logLevel string) {
	// Set the Log Level, extracted from the command line
	switch logLevel {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}
}

func printVersion() {
	fmt.Println("Dictator")
	fmt.Println("Version :",Version)
	fmt.Println("Build Time :",BuildTime)
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
	// Init flags, to get logLevel and configuration file name
	logLevel, configurationFileName := initFlags()

	setLogLevel(logLevel)
	log.Info("Dictator: Starting")
    
	// Load the configuration from config file. If something wrong, the full process is stopped inside the function
	dictatorConfiguration := initConfiguration(configurationFileName)

	// If the log level is setted inside the configuration file, override the command line level
	if dictatorConfiguration.LogLevel != "" {
		setLogLevel(dictatorConfiguration.LogLevel)
	}

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