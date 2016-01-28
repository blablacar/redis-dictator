package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"flag"
	"time"
	"os/signal"
	"fmt"
	"strings"
)

var (
	Version = "No Version Defined"
	BuildTime = "1970-01-01 00:00:00 (UTC)"
	GitRevision  = "No Git Rev Defined"
)

// Manage OS Signal, only for shutdown purpose
// When termination signal is received, we send a message to a chan
func manageSignal(c <-chan os.Signal, stop chan <-bool) {
	for {
		select {
		case _signal := <-c:
			if _signal == os.Kill {
				log.Debug("Kill Signal Received")
			}
			if _signal == os.Interrupt {
				log.Debug("Interrupt Signal Received")
			}
			log.Debug("Shutdown Signal Received")
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
	var logLevel = flag.String("log-level", "", "A value to choose between [DEBUG INFO WARN FATAL], can be overriden by config file")
	// The configuration filename
	var configurationFileName = flag.String("config", "", "the complete filename of the configuration file")
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

func setLogLevel(logLevel string) {
	lvl, err := log.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		log.WithField("level", logLevel).Fatal("Invalid log level")
	}
	log.SetLevel(lvl)
}

func printVersion() {
	fmt.Println("Version :", Version)
	fmt.Println("Build Time :", BuildTime)
	fmt.Println("Git Revision :", GitRevision)
}

func initConfiguration(configFilePath string) (DictatorConfiguration) {
	dictatorConfig := NewDictatorConfiguration()
	err := dictatorConfig.ReadConfigurationFile(configFilePath)
	if err != nil {
		log.WithError(err).Fatal("Unable to load Configuration")
	}
	return dictatorConfig
}

func main() {
	// Init flags, to get logLevel and configuration file name
	cliLogLevel, configFilePath := initFlags()

	// Load the configuration from config file. If something wrong, the full process is stopped inside the function
	dictatorConfiguration := initConfiguration(configFilePath)

	// If the log level is setted inside the configuration file, override the command line level
	if dictatorConfiguration.LogLevel != "" {
		setLogLevel(dictatorConfiguration.LogLevel)
	}

	if cliLogLevel != "" {
		setLogLevel(cliLogLevel)
	}

	log.Info("Starting")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	log.Debug("Signal channel notification setup done")

	stop := make(chan bool)
	go manageSignal(c, stop)
	log.Debug("Signal management started")

	finished := make(chan bool)
	go Run(dictatorConfiguration, stop, finished)
	log.Debug("Go routine launched")

	log.Debug("Waiting for main process to Stop")
	isFinished := <-finished
	if (isFinished) {
		log.Debug("Main routine closed correctly")
	}else {
		log.Warn("Main routine closed incorrectly")
	}
	log.Info("Shutdown")
}