package config

import (
	"testing"

	"github.com/dstoy/tunnel/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestParse(t *testing.T) {
	config, err := parse([]byte(
		"queue:\n" +
			"    url: url\n",
	))

	assert.Nil(t, err)
	assert.Equal(t, config.Queue, "url")
	assert.Equal(t, len(config.Triggers), 0)
}

func TestParseSingleTrigger(t *testing.T) {
	config, err := parse([]byte(
		"queue:\n" +
			"    url: url\n" +
			"triggers: \n" +
			"    - event: command\n",
	))

	assert.Nil(t, err)
	assert.Equal(t, config.Queue, "url")
	assert.Equal(t, len(config.Triggers), 1)

	var trigger = config.Triggers[0]
	assert.Equal(t, trigger.Event, "event")
	assert.Equal(t, trigger.Command, "command")
}

func TestParseMultipleTriggers(t *testing.T) {
	config, err := parse([]byte(
		"queue:\n" +
			"    url: url\n" +
			"triggers: \n" +
			"    - event1: command1\n" +
			"    - event2: command2\n",
	))

	assert.Nil(t, err)
	assert.Equal(t, config.Queue, "url")
	assert.Equal(t, len(config.Triggers), 2)

	var trigger = config.Triggers[0]
	assert.Equal(t, trigger.Event, "event1")
	assert.Equal(t, trigger.Command, "command1")

	trigger = config.Triggers[1]
	assert.Equal(t, trigger.Event, "event2")
	assert.Equal(t, trigger.Command, "command2")
}

func TestParseError(t *testing.T) {
	config, err := parse([]byte("queue:\n url url\n"))

	assert.NotNil(t, err)
	assert.Nil(t, config)
}

func TestGetInstance(t *testing.T) {
	var config = Instance()
	assert.Equal(t, config, instance)
}

func TestLoadDefaultConfig(t *testing.T) {
	mockReader := new(mocks.MockReader)
	mockReader.On("ReadFile", mock.Anything).Return("", nil)

	reader = mockReader

	// load default configuration
	Load()
	mockReader.AssertCalled(t, "ReadFile", CONFIG)
}

func TestLoadCustomConfig(t *testing.T) {
	mockReader := new(mocks.MockReader)
	mockReader.On("ReadFile", mock.Anything).Return("", nil)

	reader = mockReader

	// load default configuration
	Load("custom.yaml")
	mockReader.AssertCalled(t, "ReadFile", "custom.yaml")
}
