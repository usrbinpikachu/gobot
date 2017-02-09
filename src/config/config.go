package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

//Config is a struct for storing the toml data from the config file.
type Config struct {
	Server    string
	Port      int
	Channel   string
	Botname   string
	Whitelist Whitelist `toml:"Whitelist"`
}

//Whitelist is a struct for storing the whitelisted users list.
type Whitelist struct {
	Users []string
}

//ReadConfig reads in and parses toml data from the config file.
func ReadConfig() Config {
	var config Config
	if _, err := toml.DecodeFile("./gobot.conf", &config); err != nil {
		log.WithFields(log.Fields{
			"botStartup": "configLoad",
			"status": "failure",
		}).Fatal(fmt.Sprintf("Error loading config file: %s", err))
	}

	return config
}
