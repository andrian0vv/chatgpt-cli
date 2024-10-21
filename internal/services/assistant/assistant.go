package assistant

//go:generate mockgen -source $GOFILE -destination mocks/$GOFILE -package mocks

import (
	"context"
	"errors"
	"fmt"

	"github.com/andrian0vv/chatgpt-cli/internal/dto"
	"github.com/andrian0vv/chatgpt-cli/internal/logger"
)

type client interface {
	Model() string
	CreateChatCompletion(ctx context.Context, chat *dto.Chat) (string, error)
	ModelExists(ctx context.Context) (bool, error)
	GetModels(ctx context.Context) ([]string, error)
}

// Assistant is a service that provides an interface to interact with the AI assistant.
type Assistant struct {
	client client
	log    *logger.Logger
}

func New(ctx context.Context, client client, log *logger.Logger) (*Assistant, error) {
	a := &Assistant{
		client: client,
		log:    log,
	}

	if err := a.validateModel(ctx); err != nil {
		return nil, fmt.Errorf("validate model: %w", err)
	}

	return a, nil
}

// Model returns the current model name.
func (a *Assistant) Model() string {
	return a.client.Model()
}

// GetModels returns a list of available models.
func (a *Assistant) GetModels(ctx context.Context) ([]string, error) {
	m, err := a.client.GetModels(ctx)
	if err != nil {
		return nil, fmt.Errorf("get models: %w", err)
	}

	return m, nil
}

// SendMessage sends a message to the AI assistant and returns the response.
func (a *Assistant) SendMessage(ctx context.Context, question string) (string, error) {
	return a.SendChatMessage(ctx, dto.NewChat(), question)
}

// SendChatMessage sends a message to the AI assistant in the context of a chat and returns the response.
func (a *Assistant) SendChatMessage(ctx context.Context, chat *dto.Chat, question string) (string, error) {
	if question == "" {
		return "", errors.New("empty question")
	}

	if chat == nil {
		chat = dto.NewChat()
	}

	chat.AddMessage(dto.RoleUser, question)

	answer, err := a.client.CreateChatCompletion(ctx, chat)
	if err != nil {
		return "", fmt.Errorf("create chat completion: %w", err)
	}

	chat.AddMessage(dto.RoleAssistant, answer)

	return answer, nil
}

// validateModel checks if the current model exists.
func (a *Assistant) validateModel(ctx context.Context) error {
	exists, err := a.client.ModelExists(ctx)
	if err != nil {
		return fmt.Errorf("model exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("model %s does not exist", a.client.Model())
	}

	return nil
}
