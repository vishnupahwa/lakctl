package control

import (
	"context"
	"errors"
	"github.com/mitchellh/go-ps"
	"github.com/vishnupahwa/lakctl/cmd/commands/options"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// Start the command and the server
func Start(ctx context.Context, run *options.Run, serve *options.Server) error {
	cmdCtx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()
	cmd := createCommand(cmdCtx, run)
	must(cmd.Start())
	cmdPtr := &cmd
	runCtlServer(ctx, cmdPtr, cmdCtx, run, serve)
	log.Println("lakctl closed")
	return killGroupForProcess(*cmdPtr)
}

func createCommand(cmdCtx context.Context, run *options.Run) *exec.Cmd {
	c := exec.CommandContext(cmdCtx, "bash", "-c", run.Command)
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c
}

// killGroupForProcess checks the PID is running and is a child process of lakctl before killing it's group and itself.
// The cmd process kill and the wait is a safety check to make sure the process tree has fully been terminated.
func killGroupForProcess(cmd *exec.Cmd) error {
	pid, err := ps.FindProcess(cmd.Process.Pid)
	if err != nil {
		log.Printf("Cannot find PID %d: %v\n", cmd.Process.Pid, err)
		return nil
	}
	if pid == nil {
		log.Printf("Command previously running with %d is not running\n", cmd.Process.Pid)
		return nil
	}
	if pid.PPid() != os.Getpid() {
		log.Printf("Process %d has parent %d and is not a subprocess of lakctl (%d).", pid.Pid(), pid.PPid(), os.Getpid())
		return nil
	}
	log.Printf("Stopping %s (PID: %d, PPID: %d)", pid.Executable(), pid.Pid(), pid.PPid())
	errGroup := syscall.Kill(-pid.Pid(), syscall.SIGTERM)
	errCmd := cmd.Process.Signal(syscall.SIGTERM)
	waitWithTimeout(cmd)
	return compositeErr(errCmd, errGroup)
}

func waitWithTimeout(cmd *exec.Cmd) {
	wait := make(chan error, 1)
	go func() { wait <- cmd.Wait() }()
	select {
	case <-wait:
		return
	case <-time.After(30 * time.Second):
		log.Printf("Timed out waiting 30s for process to be killed\n")
	}
}

func compositeErr(errs ...error) error {
	message := ""
	for _, err := range errs {
		if err != nil {
			message = message + err.Error() + " "
		}
	}
	if len(message) > 0 {
		return errors.New(message)
	}
	return nil
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
