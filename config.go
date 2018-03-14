package main

import (
	"log"
	"os"

	"github.com/junwangustc/ThirdAlarm/alert"
	"github.com/junwangustc/ThirdAlarm/lark"
	"github.com/junwangustc/ThirdAlarm/slack"
	"github.com/junwangustc/ThirdAlarm/sms"
	"github.com/naoina/toml"
)

type Config struct {
	Slack slack.Config `toml:"slack", json:"slack"`
	Alert alert.Config `toml:"alert", json:"alert"`
	SMS   sms.Config   `toml:"sms"`
	Lark  lark.Config  `toml:"lark"`
}

func NewConfig() *Config {
	c := &Config{}
	c.Alert = alert.NewConfig()
	c.Slack = slack.NewConfig()
	c.SMS = sms.NewConfig()
	c.Lark = lark.NewConfig()
	return c
}

func ParseConfig(path string) (cfg *Config, err error) {
	if path == "" {
		log.Fatalln("no configuration provided, using default settings")
	}
	log.Printf("Using configuration at: %s\n", path)
	config := NewConfig()
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close()
	return config, toml.NewDecoder(f).Decode(&config)
}
