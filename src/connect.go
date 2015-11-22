package gobot

import (
	"fmt"
	"strconv"
	"strings"

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
