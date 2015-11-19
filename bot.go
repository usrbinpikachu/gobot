package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/thoj/go-ircevent"
)

type Config struct {
	Server  string
	Port    int
	Channel string
	Botname string
}

func ReadConfig() Config {
	var config Config
	if _, err := toml.DecodeFile("./gobot.conf", &config); err != nil {
		fmt.Println(err)
	}

	return config
}

func main() {
	var config = ReadConfig()
	var channel = config.Channel

	//The IRC function takes a user and nick, we just send the same thing for both.
	connection := irc.IRC(config.Botname, config.Botname)

	//Override the irc-event's logging to stdout to log to a file.
	logFile, loggerErr := os.OpenFile("gobot.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if loggerErr != nil {
		fmt.Printf("Error opening log file: %v", loggerErr)
	}
	defer logFile.Close()
	connection.Log = log.New(logFile, "", log.LstdFlags)

	connectionString := []string{config.Server, strconv.Itoa(config.Port)}
	err := connection.Connect(strings.Join(connectionString, ":"))
	if err != nil {
		fmt.Printf("Connection error: %s", err)
		return
	}

	connection.AddCallback("001", func(e *irc.Event) {
		connection.Join(channel)
	})

	connection.Loop()
}
