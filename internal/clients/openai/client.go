package openai

import (
	"context"
	"fmt"
	"sort"

	"github.com/sashabaranov/go-openai"

	"github.com/andrian0vv/chatgpt-cli/internal/config"
	"github.com/andrian0vv/chatgpt-cli/internal/dto"
	"github.com/andrian0vv/chatgpt-cli/internal/logger"
)

const defaultModel = openai.GPT3Dot5Turbo

type Client struct {
	client *openai.Client
	model  string
	log    *logger.Logger
}

func New(cfg config.Config, log *logger.Logger, opts ...Option) *Client {
	c := &Client{
		client: openai.NewClient(cfg.OpenaiApiKey),
		model:  defaultModel,
		log:    log,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) Model() string {
	return c.model
}

func (c *Client) CreateChatCompletion(ctx context.Context, chat *dto.Chat) (string, error) {
	in := c.toCreateChatCompletionIn(chat)

	c.log.Debug("openai in CreateChatCompletion", logger.WithField("in", in))

	out, err := c.client.CreateChatCompletion(ctx, in)
	if err != nil {
		return "", fmt.Errorf("create chat completion: %w", err)
	}

	c.log.Debug("openai out CreateChatCompletion", logger.WithField("out", out))

	if len(out.Choices) == 0 {
		return "", fmt.Errorf("empty answer")
	}

	return out.Choices[0].Message.Content, nil
}

func (c *Client) ModelExists(ctx context.Context) (bool, error) {
	if c.model == defaultModel {
		return true, nil
	}

	out, err := c.client.GetModel(ctx, c.model)
	if err != nil {
		return false, fmt.Errorf("get model: %w", err)
	}

	c.log.Debug("openai out GetModel", logger.WithField("out", out))

	return out.ID == c.model, nil
}

func (c *Client) GetModels(ctx context.Context) ([]string, error) {
	out, err := c.client.ListModels(ctx)
	if err != nil {
		return nil, fmt.Errorf("list models: %w", err)
	}

	c.log.Debug("openai out ListModels", logger.WithField("out", out))

	list := make([]string, 0, len(out.Models))
	for _, model := range out.Models {
		list = append(list, model.ID)
	}

	sort.Strings(list)

	c.log.Debug("models", logger.WithField("models", list))

	return list, nil
}
