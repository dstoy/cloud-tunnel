package command

import (
	"log"
	"os/exec"
)

type Logger struct {
}

type Command struct {
	Run  string
	User string
}

var logger = &Logger{}

/**
* Start executing a command
 */
func Dispatch(command *Command) error {
	log.Println("Executing", command.Run)

	var cmd *exec.Cmd
	if command.User != "" {
		cmd = exec.Command("sudo", "-u", command.User, "bash", "-c", command.Run)
	} else {
		cmd = exec.Command("bash", "-c", command.Run)
	}

	cmd.Stdout = logger
	cmd.Stderr = logger

	cmd.Run()
	log.Println("Execution complete")

	return nil
}

/**
* Implement the Write method for the logger interface
 */
func (so *Logger) Write(msg []byte) (n int, err error) {
	log.Printf("%s", msg)
	return len(msg), nil
}
