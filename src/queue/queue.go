package queue

import (
	"log"

	"github.com/dstoy/tunnel/src/command"
	"github.com/dstoy/tunnel/src/config"
)

type Message struct {
	Event string
	Id    string
}

type Queue interface {
	Connect(*string) error
	GetMessage() *Message
	DeleteMessage(*Message) error
}

var queue Queue = &SQS{}
var conf = config.Instance()

/**
* Listen to the queue
 */
func Listen() error {
	err := queue.Connect(&conf.Queue)
	if err != nil {
		return err
	}

	// Listen for messages
	for {
		message := queue.GetMessage()
		cmd := getCommand(message.Event)
		log.Println("Received event:", message.Event)

		if cmd != "" {
			command.Dispatch(cmd)
		} else {
			log.Println(
				"Received invalid event without a matching action:",
				message.Event,
			)
		}

		err := queue.DeleteMessage(message)
		if err != nil {
			log.Println("Error removing a message", err)
		}
	}
}

/**
* Return the command for a supplied event
 */
func getCommand(event string) string {
	for i := 0; i < len(conf.Triggers); i++ {
		if conf.Triggers[i].Event == event {
			return conf.Triggers[i].Command
		}
	}

	return ""
}
