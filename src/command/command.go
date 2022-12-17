package command

import (
	"log"
	"os/exec"
)

type Logger struct {
}

var logger = &Logger{}

/**
* Start executing a command
 */
func Dispatch(action string) error {
	log.Println("Executing", action)

	cmd := exec.Command("bash", "-c", action)
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
