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
		Long: `Start an application and an HTTP server in order to control it.

Usage:
$ lakctl start --run "go run main.go Hello World"

Stopping the running application:
$ curl http://localhost:8008/stop

Starting the application:
$ curl http://localhost:8008/start

Restarting the application substituting the --run flag: (/restart/<old base64 encoded value>/<new base64 encoded value>)
$ curl http://localhost:8008/restart/V29ybGQK/5LiW55WMCg==

N.B. The start and restart endpoints will always try and start or substitute the command given in the '--run' flag.
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := sigContext()
			err := control.Start(ctx, runOpts, serverOpts)
			if err != nil {
				log.Fatal(err)
			}
		},
		Aliases: []string{"init", "run"},
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
