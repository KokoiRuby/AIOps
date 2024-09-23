package main

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	token := "sk-NUkr3DrnPJtiWZS32d667dAaCd304f20B43bAbF0D6872b21"
	config := openai.DefaultConfig(token)
	config.BaseURL = "https://vip.apiyi.com/v1"
	client := openai.NewClientWithConfig(config)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: "json_object",
			},
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a JSON parser. Please output a json object based on example. Demo: {\"service_name\":\"\",\"action\":\"\"}，where action could be get_log, restart, delete",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Please help me restart payment service.",
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
