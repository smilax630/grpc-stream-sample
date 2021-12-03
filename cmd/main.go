package main

import (
	"os"

	"github.com/grpc-streamer/pkg/interface/cmd"
	"github.com/spf13/cobra"
)

func main() {
	c := &cobra.Command{Use: "gateway [command]"}
	cmd.RegisterCommand(c)
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
