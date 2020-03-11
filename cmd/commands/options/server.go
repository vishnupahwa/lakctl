package options

import (
	"github.com/spf13/cobra"
	"strconv"
)

// Server struct contains server options for the controlling HTTP server.
type Server struct {
	Port int
}

func AddServerArgs(cmd *cobra.Command, s *Server) {
	cmd.Flags().IntVarP(&s.Port, "port", "p", 8008,
		"Port to use for HTTP control server")
}

func (s *Server) PortStr() string {
	return strconv.Itoa(s.Port)
}
