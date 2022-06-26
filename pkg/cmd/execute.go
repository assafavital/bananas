package cmd

import (
	"context"

	"github.com/spf13/cobra"
	cf "raftt.io/bananas/pkg/controlflow"
	l "raftt.io/bananas/pkg/logging"
)

type Command interface {
	Execute()
}

func execute(ctx context.Context, command *cobra.Command) error {
	return cf.Execute(func() error { return command.ExecuteContext(ctx) })
}

func wrapExecute(ctx context.Context, command *cobra.Command) {
	if err := execute(ctx, command); err != nil {
		l.Logger.WithError(err).Error("Failed executing command")
		l.Logger.Exit(1)
	}
}
