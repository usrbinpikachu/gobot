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

//Config structure for storing the toml data from the config file.
type Config struct {
	Server  string
	Port    int
	Channel string
	Botname string
}

//ReadConfig reads in and parses toml data from the config file.
func ReadConfig() Config {
	var config Config
	if _, err := toml.DecodeFile("./gobot.conf", &config); err != nil {
		fmt.Println(err)
	}

	return config
}

//Connect establishes a connection to the IRC server specified in gobot.conf.
func Connect(botName string, botUsername string, serverAddress string, serverPort int) *irc.Connection {
	connection := irc.IRC(botName, botUsername)

	//TODO: There's probably a better way to do this than strings.Join().
	connectionString := []string{serverAddress, strconv.Itoa(serverPort)}
	err := connection.Connect(strings.Join(connectionString, ":"))
	if err != nil {
		fmt.Printf("Connection error: %s", err)
		return nil
	}

	return connection
}

func main() {
	config := ReadConfig()
	channel := config.Channel

	//The IRC function takes a nick and username, we just send the same thing for both.
	connection := Connect(config.Botname, config.Botname, config.Server, config.Port)

	//Override irc-event's default logging to stdout to log to a file.
	logFile, loggerErr := os.OpenFile("gobot.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if loggerErr != nil {
		fmt.Printf("Error opening log file: %v", loggerErr)
	}
	defer logFile.Close()
	connection.Log = log.New(logFile, "", log.LstdFlags)

	//001 is the WELCOME event, which means we successfully connected.
	connection.AddCallback("001", func(e *irc.Event) {
		connection.Join(channel)
	})

	connection.AddCallback("PRIVMSG", func(e *irc.Event) {
		connection.Log.Printf(e.Message())
	})

	connection.Loop()
}
