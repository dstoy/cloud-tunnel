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
	Connect(*config.QueueConfig) error
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

		if cmd != nil {
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
func getCommand(event string) *command.Command {
	for i := 0; i < len(conf.Triggers); i++ {
		trigger := conf.Triggers[i]
		if trigger.Event == event {
			return &command.Command{
				Run:  trigger.Run,
				User: trigger.User,
			}
		}
	}

	return nil
}
