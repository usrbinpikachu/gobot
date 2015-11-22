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

//CheckWhitelist checks incoming events' sender nicks against the whitelist.
func CheckWhitelist(e *irc.Event, c Config) bool {
	for _, s := range c.Whitelist.Users {
		if s == e.Nick {
			return true
		}
	}

	return false
}

func main() {
	config := ReadConfig()
	channel := config.Channel

	//The IRC function takes a nick and username, we send the same thing for both.
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

	//On PRIVMSG log the nick and message, then check if the nick is whitelisted.
	connection.AddCallback("PRIVMSG", func(e *irc.Event) {
		connection.Log.Printf("%s: %s", e.Nick, e.Message())
		if CheckWhitelist(e, config) {
			connection.Log.Printf("%s is whitelisted.", e.Nick)
		} else {
			connection.Log.Printf("%s is not whitelisted.", e.Nick)
		}
	})

	connection.Loop()
}
