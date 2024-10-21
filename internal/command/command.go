package command

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/fatih/color"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"

	"github.com/andrian0vv/chatgpt-cli/internal/clients/openai"
	"github.com/andrian0vv/chatgpt-cli/internal/config"
	"github.com/andrian0vv/chatgpt-cli/internal/logger"
	"github.com/andrian0vv/chatgpt-cli/internal/services/assistant"
)

const (
	colorSystem = color.FgHiBlue
	colorAI     = color.FgYellow
	colorError  = color.FgRed
)

// Command is a wrapper around cobra.Command with additional printing methods.
type Command struct {
	*cobra.Command
	Assistant *assistant.Assistant
}

func New(c *cobra.Command) Command {
	cmd := Command{
		Command: c,
	}

	cmd.Assistant = cmd.createAssistant()

	return cmd
}

func (c Command) createAssistant() *assistant.Assistant {
	cfg := config.New()

	verbose, err := c.Flags().GetBool("verbose")
	c.Fail(err)

	model, err := c.Flags().GetString("model")
	c.Fail(err)

	log := logger.New(c.OutOrStdout(), logger.WithEnabled(verbose))

	a, err := assistant.New(
		c.Context(),
		openai.New(cfg, log, openai.WithModel(model)),
		log,
	)
	c.Fail(err)

	return a
}

func (c Command) Clear() {
	c.Print("\r\033[K")
}

func (c Command) Fail(err error) {
	if err == nil {
		return
	}

	c.print(colorError, "[Error] %v\n", err)

	os.Exit(1)
}

func (c Command) AI(message string) {
	if strings.Contains(message, "\n") {
		r, err := glamour.NewTermRenderer(
			glamour.WithColorProfile(termenv.ANSI256),
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(100),
		)
		c.Fail(err)

		message, err = r.Render(message)
		c.Fail(err)

		message = strings.Trim(message, "\n")
	}

	c.print(colorAI, "[AI] %s\n", message)
}

func (c Command) System(message string) {
	c.print(colorSystem, "[System] %s\n", message)
}

func (c Command) Loading(message string) func() {
	done := make(chan struct{})

	go func() {
		const maxDotsCount = 3
		dotsCount := maxDotsCount

		printDotsFn := func(counter int) {
			c.print(colorAI, "[AI] %s%s", message, strings.Repeat(".", counter))
		}

		printDotsFn(dotsCount)

		for {
			if dotsCount > maxDotsCount {
				dotsCount = 1
			}

			select {
			case <-done:
				return
			case <-time.After(400 * time.Millisecond):
				c.Clear()
				printDotsFn(dotsCount)
				dotsCount++
			}
		}
	}()

	return func() {
		done <- struct{}{}
		close(done)
		c.Clear()
	}
}

func (c Command) Read() string {
	c.Print("[You] ")

	question, _ := bufio.NewReader(c.InOrStdin()).ReadString('\n')

	return strings.TrimSpace(question)
}

func (c Command) print(attribute color.Attribute, message string, args ...any) {
	_, _ = color.New(attribute).Fprintf(c.OutOrStdout(), message, args...)
}
