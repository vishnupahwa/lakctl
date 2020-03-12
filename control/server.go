package control

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/vishnupahwa/lakctl/cmd/commands/options"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

// runCtlServer starts the HTTP server for controlling the command.
func runCtlServer(ctx context.Context, cmd **exec.Cmd, cmdCtx context.Context, run *options.Run, serve *options.Server) {
	http.HandleFunc("/start", handleStart(cmd, cmdCtx, run))
	http.HandleFunc("/stop", handleStop(cmd))
	http.HandleFunc("/restart/", handleRestart(cmd, cmdCtx, run))
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

func handleRestart(cmd **exec.Cmd, ctx context.Context, run *options.Run) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		split := strings.Split(r.URL.Path, "/")
		if len(split) < 4 {
			fmt.Printf("Expected path of /restart/<base64 old string>/<base64 new string> but got %s\n", r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dec := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(split[2]))
		res, err := ioutil.ReadAll(dec)
		oldStr := strings.TrimSuffix(string(res), "\n")
		dec = base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(split[3]))
		res, err = ioutil.ReadAll(dec)
		if err != nil {
			fmt.Printf("Failed to base64 decode %s and %s: %v\n", split[2], split[3], err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		newStr := strings.TrimSuffix(string(res), "\n")

		err = killGroupForProcess(*cmd)
		if err != nil {
			log.Printf("Failed to kill process %v (continuing)\n", err)
		}

		log.Printf(`Replacing %s with %s in command: "%s"`, oldStr, newStr, run.Command)
		subRun := options.Run{
			Command: strings.ReplaceAll(run.Command, oldStr, newStr),
		}
		newCommand := createCommand(ctx, &subRun)
		must(newCommand.Start())
		log.Printf("Started substituted command with PID %d\n", newCommand.Process.Pid)
		*cmd = newCommand
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
