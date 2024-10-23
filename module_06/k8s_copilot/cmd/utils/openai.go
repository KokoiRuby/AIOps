package utils

import (
	"context"
	"errors"
	"github.com/sashabaranov/go-openai"
	"os"
)

type OpenAI struct {
	Client *openai.Client
	ctx    context.Context
}

func NewOpenAI() (*OpenAI, error) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, errors.New("API_KEY environment variable is not set")
	}
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		return nil, errors.New("BASE_URL environment variable is not set")
	}
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL
	client := openai.NewClientWithConfig(config)

	ctx := context.Background()

	return &OpenAI{
		Client: client,
		ctx:    ctx,
	}, nil
}

func (o *OpenAI) SendMessage(prompt, input string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: prompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: input,
			},
		},
	}
	resp, err := o.Client.CreateChatCompletion(o.ctx, req)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", errors.New("no choices found")
	}
	return resp.Choices[0].Message.Content, nil
}
