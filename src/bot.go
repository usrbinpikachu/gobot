package main

import (
	"fmt"
	"os"

	"github.com/thoj/go-ircevent"
	log "github.com/sirupsen/logrus"

	"./config"
	"./connect"
	"./dictionary"
	"./wunderground"
)

//CheckWhitelist checks incoming events' sender nicks against the whitelist.
func CheckWhitelist(e *irc.Event, c config.Config) bool {
	for _, s := range c.Whitelist.Users {
		if s == e.Nick {
			return true
		}
	}

	return false
}

func main() {
	readConfig := config.ReadConfig()
	if &readConfig == nil {
		log.WithFields(log.Fields{
			"botStartup": "configLoad",
			"status": "failure",
		}).Fatal("Config not loaded. Is the path correct?")
	}
	channel := readConfig.Channel

	//Initialize logrus logger.
	log.SetFormatter(&log.JSONFormatter{})

	//Override irc-event's default logging to stdout to log to a file.
	logFile, loggerErr := os.OpenFile("gobot.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if loggerErr != nil {
		log.WithFields(log.Fields{
			"botStartup": "logInit",
			"status": "failure",
		}).Error("Failed open log file, using default stderr.")
	} else {
		log.SetOutput(logFile)
	}

	defer logFile.Close()

	//The IRC function takes a nick and username, we send the same thing for both.
	connection := connect.Connect(readConfig.Botname, readConfig.Botname, readConfig.Server, readConfig.Port)


	//001 is the WELCOME event, which means we successfully connected.
	connection.AddCallback("001", func(e *irc.Event) {
		connection.Join(channel)
		log.WithFields(log.Fields{
			"event": "channelJoin",
			"status": "success",
		}).Info("Successfully joined %s.", channel)
	})

	//On PRIVMSG log the nick and message, then check if the nick is whitelisted.
	connection.AddCallback("PRIVMSG", func(e *irc.Event) {
		log.WithFields(log.Fields{
			"event": "PRIVMSG",
			"sender": e.Nick,
		}).Info(e.Message())

		if CheckWhitelist(e, readConfig) {
			log.WithFields(log.Fields{
				"event": "PRIVMSG",
				"sender": e.Nick,
			}).Info(fmt.Sprintf("%s is whitelisted.", e.Nick))
		} else {
			log.WithFields(log.Fields{
				"event": "PRIVMSG",
				"sender": e.Nick,
			}).Info(fmt.Sprintf("%s is not whitelisted.", e.Nick))
		}
	})

	temp, err := wunderground.Temperature("98004")
	if err != nil {
		log.WithFields(log.Fields{
			"event": "wunderground",
			"status": "failure",
		}).Error(fmt.Sprintf("Error retrieving Wunderground API data: %s", err))

	}
	log.WithFields(log.Fields{
		"event": "wunderground",
		"status": "success",
	}).Info(fmt.Sprintf("Temperature %gF", temp))

	//TODO: Rework this to return all of the returned definitions instead of just the first.
	word, err := dictionary.Define("cake")
	if err != nil {
		log.WithFields(log.Fields{
			"event": "dictionary",
			"status": "failure",
		}).Error(fmt.Sprintf("Error retrieving Wordnik API data: %s", err))
	}
	log.WithFields(log.Fields{
		"event": "dictionary",
		"status": "success",
	}).Info(fmt.Sprintf("%s: %s", word[0].Word, word[0].Definition))

	connection.Loop()
}
