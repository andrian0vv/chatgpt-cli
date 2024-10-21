package ask

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/andrian0vv/chatgpt-cli/internal/command"
)

var Command = &cobra.Command{
	Use:   "ask",
	Short: "Ask AI with one question",
	Args:  cobra.MatchAll(cobra.MinimumNArgs(1)),
	Run:   Run,
}

func Run(c *cobra.Command, args []string) {
	cmd := command.New(c)

	question := strings.Join(args, " ")

	answer, err := cmd.Assistant.SendMessage(cmd.Context(), question)
	cmd.Fail(err)

	cmd.AI(answer)
}
