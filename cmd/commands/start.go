package commands

import (
	"context"
	"github.com/vishnupahwa/lakctl/cmd/commands/options"
	"github.com/vishnupahwa/lakctl/control"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

// addStart adds the primary start command to a top level command. This is the entrypoint command for starting
// a controlled application.
func addStart(topLevel *cobra.Command) {
	serverOpts := &options.Server{}
	runOpts := &options.Run{}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start an application and an HTTP server to control it",
		Long: `Start an application and an HTTP server for starting, stopping and restarting with different commands.
N.B. All HTTP endpoints which start the command use the '--run' command. 
lakctl does not continuously store modified commands, only the running command's PID.
When passing a substitution to /restart multiple times, each substitution will replace the original input
to the --run command.
e.g (the restart inputs are base64 decoded for clarity)
$ lakctl start --run "go run main.go hello world"
$ curl localhost:8008/restart/world/galaxy 
$ curl localhost:8008/restart/world/universe 

The final command would cause "go run main.go hello universe" to be run, even if the previous command
was ...hello galaxy.

`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := sigContext()
			err := control.Start(ctx, runOpts, serverOpts)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	options.AddServerArgs(startCmd, serverOpts)
	options.AddRunArg(startCmd, runOpts)
	topLevel.AddCommand(startCmd)
}

// sigContext creates a context which will be cancelled on a SIGINT or SIGTERM
func sigContext() context.Context {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancelCtx := context.WithCancel(context.Background())
	go cancelIfNotified(signals, cancelCtx)
	return ctx
}

func cancelIfNotified(signals chan os.Signal, cancelCtx context.CancelFunc) {
	<-signals
	cancelCtx()
}
