package openai

import (
	"github.com/sashabaranov/go-openai"

	"github.com/andrian0vv/chatgpt-cli/internal/dto"
)

const maxMessages = 20

func (c *Client) toCreateChatCompletionIn(chat *dto.Chat) openai.ChatCompletionRequest {
	chatMessages := chat.Messages
	if len(chatMessages) > maxMessages {
		chatMessages = chatMessages[len(chatMessages)-maxMessages:]
	}

	messages := make([]openai.ChatCompletionMessage, 0, len(chatMessages))
	for _, message := range chatMessages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    string(message.Role),
			Content: message.Content,
		})
	}

	return openai.ChatCompletionRequest{
		Model:       c.model,
		Messages:    messages,
		Temperature: 0.7,
		N:           1,
	}
}
