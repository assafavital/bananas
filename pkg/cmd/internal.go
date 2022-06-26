package cmd

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"raftt.io/bananas/pkg/config"
	l "raftt.io/bananas/pkg/logging"
)

type internalCommand struct {
	Command *cobra.Command
}

func (cmd internalCommand) Execute() {
	wrapExecute(context.Background(), cmd.Command)
}

func MakeInternalCommand(command *cobra.Command, logfilePath string) Command {
	logfilePath = config.LogFilePath()
	l.Setup(logrus.TraceLevel, logfilePath)
	return internalCommand{Command: command}
}
