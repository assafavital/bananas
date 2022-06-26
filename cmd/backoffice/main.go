package main

import (
	"raftt.io/bananas/cmd/backoffice/app"
	"raftt.io/bananas/pkg/cmd"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:  "server",
		Long: "Raftt backoffice server",
	}

	runCmd = &cobra.Command{
		Use:  "run",
		Long: "Run backoffice server handling user authentication and Github OAuth",
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.RunServer(cmd.Context())
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
}

func main() {
	cmd.MakeInternalCommand(rootCmd, "").Execute()
}
