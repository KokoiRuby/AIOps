package main

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	token := "BLOCKED"
	config := openai.DefaultConfig(token)
	config.BaseURL = "https://vip.apiyi.com/v1"
	client := openai.NewClientWithConfig(config)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You're an operation & maintenance expert. Your job is to help user solve technical problems.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Linux ports are in conflict, how to troubleshoot? Just give me a few commands.",
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
