package models

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/andrian0vv/chatgpt-cli/internal/command"
)

var Command = &cobra.Command{
	Use:   "models",
	Short: "Get list of ChatGPT models",
	Args:  cobra.MatchAll(cobra.NoArgs),
	Run:   Run,
}

func Run(c *cobra.Command, _ []string) {
	cmd := command.New(c)

	models, err := cmd.Assistant.GetModels(cmd.Context())
	cmd.Fail(err)

	var answer strings.Builder
	for _, model := range models {
		answer.WriteString(fmt.Sprintf("* %s\n", model))
	}

	cmd.AI(answer.String())
}
