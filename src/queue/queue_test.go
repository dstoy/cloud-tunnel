package queue

import (
	"testing"

	"github.com/dstoy/tunnel/src/config"
	"github.com/stretchr/testify/assert"
)

func TestGetCommand(t *testing.T) {
	conf = &config.Config{
		Queue: config.QueueConfig{
			Url: "test",
		},
		Triggers: []config.Trigger{
			{Event: "event", Run: "command"},
		},
	}

	command := getCommand("event")
	assert.Equal(t, command, "command")

	command = getCommand("other")
	assert.Equal(t, command, "")
}
