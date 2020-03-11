package control

import (
	"context"
	"fmt"
	"github.com/vishnupahwa/lakctl/cmd/commands/options"
	"log"
	"net/http"
	"os/exec"
)

// runCtlServer starts the HTTP server for controlling the command.
func runCtlServer(ctx context.Context, cmd **exec.Cmd, cmdCtx context.Context, run *options.Run, serve *options.Server) {
	http.HandleFunc("/start", handleStart(cmd, cmdCtx, run))
	http.HandleFunc("/stop", handleStop(cmd))
	http.HandleFunc("/", handle(cmd))
	server := &http.Server{Addr: ":" + serve.PortStr()}
	startServer(server)
	waitForServer(ctx, server)
}

func handleStart(cmd **exec.Cmd, ctx context.Context, run *options.Run) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c := *cmd
		if c.ProcessState != nil {
			newCommand := createCommand(ctx, run)
			must(newCommand.Start())
			log.Printf("Started command with PID %d\n", newCommand.Process.Pid)
			*cmd = newCommand
		} else {
			log.Printf("Command is already running with PID %d\n", c.Process.Pid)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handleStop(cmd **exec.Cmd) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := killGroupForProcess(*cmd)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func startServer(server *http.Server) {
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failure: %v", err)
		}
	}()
}

func waitForServer(ctx context.Context, server *http.Server) {
	<-ctx.Done()
	log.Println("Shutting down server")
	_ = server.Shutdown(ctx)
}

func handle(cmd **exec.Cmd) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "%d\n", (*cmd).Process.Pid)
	}
}
