package main

import (
	"fmt"
	"os"
	//"time"
)

// LiveLoggerMessage is a message for a live logger update
type LiveLoggerMessage struct {
	source    string
	eventType string
	message   string
	progress  float32
}

// LiveLogger is a the structure for live logger
type LiveLogger struct {
	msgChannel chan LiveLoggerMessage       // message channel
	open       bool                         // LiveLogger status
	frequency  int                          // polling frequency
	messages   map[string]LiveLoggerMessage // message stack to keep track of the last sent message
}

// NewLiveLogger is a function to start a new live logger
func NewLiveLogger() LiveLogger {
	ll := LiveLogger{
		msgChannel: make(chan LiveLoggerMessage),
		open:       true,
		frequency:  1000,
		messages:   make(map[string]LiveLoggerMessage),
	}
	go ll.start()
	return ll
}

func (ll LiveLogger) start() {
	for msg := range ll.msgChannel {
		ll.messages[msg.source] = msg // update the latest message
		ll.log()                      // log the updates
	}
}

func (ll LiveLogger) log() {
	fmt.Print("\033[s")                     // save cursor position
	for i := 0; i < len(ll.messages); i++ { // delete all the rows
		fmt.Print("\033[K\033[1B")
	}
	fmt.Print("\033[u") // return to saved cursor position
	fmt.Print("\033[s") // save cursor position
	for _, message := range ll.messages {
		logLine(message)
	}
	fmt.Print("\033[u")
	// fmt.Printf("\033[%dA\033[")
}

// logLine is used to log a line
func logLine(message LiveLoggerMessage) {
	var line string
	switch message.eventType {
	case "NEW":
		line = fmt.Sprintf("[ NEW ]    - %s : %s", message.source, message.message)
	case "DROP":
		line = fmt.Sprintf("[ DROP ]   - %s : %s", message.source, message.message)
	case "START":
		line = fmt.Sprintf("[ START ]  - %s : %s", message.source, message.message)
	case "UPDATE":
		line = fmt.Sprintf("[ UPDATE ] - %s ( %6.2f %% ) : %s", message.source, message.progress, message.message)
	case "END":
		line = fmt.Sprintf("[ DONE ]   - %s : %s", message.source, message.message)
	default:
		line = fmt.Sprintf("%s : %s", message.source, message.message)
	}
	fmt.Fprintln(os.Stdout, line)
}

// End is to end a live logger
// it closes the msgChannel and sets the open parameter to false
func (ll LiveLogger) End() {
	close(ll.msgChannel)
	ll.open = false
}
