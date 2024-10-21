package chat

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/andrian0vv/chatgpt-cli/cmd/models"
	"github.com/andrian0vv/chatgpt-cli/internal/command"
	"github.com/andrian0vv/chatgpt-cli/internal/dto"
)

const (
	messageOnStart = `Welcome to the dialog with ChatGPT (%s)! 
1. Type 'model' to get current model. 
2. Type 'models' to list models. 
3. Type 'reset' to reset the chat. 
4. Type 'exit' or cmd+C to stop.
`
	messageOnExit    = "Goodbye!"
	messageOnReset   = "The chat has been reset."
	messageOnLoading = "Thinking"
)

const (
	commandModel  = "model"
	commandModels = "models"
	commandReset  = "reset"
	commandExit   = "exit"
)

var Command = &cobra.Command{
	Use:   "chat",
	Short: "Start chat with AI",
	Args:  cobra.MatchAll(cobra.MaximumNArgs(1)),
	Run:   Run,
}

func Run(c *cobra.Command, _ []string) {
	cmd := command.New(c)

	cmd.System(fmt.Sprintf(messageOnStart, cmd.Assistant.Model()))

	chat := dto.NewChat()

	for {
		question := cmd.Read()

		switch question {
		case commandModel:
			cmd.AI(cmd.Assistant.Model())
		case commandModels:
			models.Run(c, nil)
		case commandReset:
			chat.Reset()
			cmd.System(messageOnReset)
		case commandExit:
			cmd.System(messageOnExit)
			return
		case "":
		default:
			cancel := cmd.Loading(messageOnLoading)

			answer, err := cmd.Assistant.SendChatMessage(cmd.Context(), chat, question)
			cancel()

			cmd.Fail(err)
			cmd.AI(answer)
		}
	}
}
