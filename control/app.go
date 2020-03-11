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
)

func Start(ctx context.Context, run *options.Run, serve *options.Server) error {
	cmdCtx, cancelFunc := context.WithCancel(ctx)
	c := createCommand(cmdCtx, run)
	c.Start()
	runCtlServer(ctx, c, cmdCtx, cancelFunc, run, serve)
	log.Println("lakctl closed")
	must(killGroupForProcess(c))
	err := c.Wait()
	return err
}

func createCommand(cmdCtx context.Context, run *options.Run) *exec.Cmd {
	c := exec.CommandContext(cmdCtx, "bash", "-c", run.Command)
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c
}

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
	errCmd := cmd.Process.Kill()
	_ = cmd.Wait()
	return compositeErr(errCmd, errGroup)
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
