package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/natefinch/lumberjack.v2"

	log "github.com/junwangustc/ustclog"
)

var (
	configFilePath string
	logFilePath    string
	globalLog      *log.Logger
)

func init() {
	flag.StringVar(&configFilePath, "config", "/tmp/config.toml", "config file")
	flag.StringVar(&logFilePath, "logfile", "/tmp/third-alarm.log", "log file")
	flag.Parse()
}

func initLog() {
	output := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     7,
	}
	log.SetOutput(output)
	globalLog = log.New(output, "", log.Ldefault)

}

func main() {
	initLog()
	cfg, err := ParseConfig(configFilePath)
	if err != nil {
		log.Fatalf("parseconfig error  start: %v", err)
	}
	srv := NewServer(cfg)
	go func() {
		if err := srv.Run(); err != nil {
			log.Fatalf("alarm start: %v", err)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	log.Printf("listening for signals ...")
	select {
	case <-signalCh:
		log.Printf("signal received, shutdown...")
	}
}
