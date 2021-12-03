package cmd

import (
	"github.com/spf13/cobra"
	// Automatically set GOMAXPROCS to match Linux container CPU quota.
	_ "go.uber.org/automaxprocs"
)

func RegisterCommand(registry *cobra.Command) {
	registry.AddCommand(
		newV1(),
	)
}
