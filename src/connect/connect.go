package connect

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/thoj/go-ircevent"
	log "github.com/sirupsen/logrus"
)

//Connect establishes a connection to the IRC server specified in gobot.conf.
func Connect(botName string, botUsername string, serverAddress string, serverPort int) *irc.Connection {
	connection := irc.IRC(botName, botUsername)

	//TODO: There's probably a better way to do this than strings.Join().
	connectionString := []string{serverAddress, strconv.Itoa(serverPort)}
	err := connection.Connect(strings.Join(connectionString, ":"))
	if err != nil {
		log.WithFields(log.Fields{
			"event": "serverConnection",
			"status": "failure",
		}).Error(fmt.Sprintf("Connection error: %s", err))
		return nil
	}

	log.WithFields(log.Fields{
		"event": "serverConnection",
		"status": "success",
	}).Info(fmt.Sprintf("Successfully connected to %s:%s", serverAddress, serverPort))
	return connection
}
