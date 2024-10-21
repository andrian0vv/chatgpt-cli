package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/andrian0vv/chatgpt-cli/cmd/ask"
	"github.com/andrian0vv/chatgpt-cli/cmd/chat"
	"github.com/andrian0vv/chatgpt-cli/cmd/models"
)

var (
	verbose bool
	model   string
)

var rootCommand = &cobra.Command{
	Use:   "chatgpt-cli",
	Short: "Smooth interaction with chat gpt",
}

func init() {
	// Commands
	rootCommand.AddCommand(ask.Command)
	rootCommand.AddCommand(chat.Command)
	rootCommand.AddCommand(models.Command)

	// Flags
	rootCommand.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCommand.PersistentFlags().StringVarP(&model, "model", "m", "", "ChatGPT model")
}

func Execute(ctx context.Context) error {
	return rootCommand.ExecuteContext(ctx)
}
