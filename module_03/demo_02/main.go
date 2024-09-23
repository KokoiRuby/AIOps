package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func main() {
	token := "sk-NUkr3DrnPJtiWZS32d667dAaCd304f20B43bAbF0D6872b21"
	config := openai.DefaultConfig(token)
	config.BaseURL = "https://vip.apiyi.com/v1"
	client := openai.NewClientWithConfig(config)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You're an operation & maintenance expert. Your job is to help user solve technical problems.",
		},
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("-> ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.Replace(userInput, "\n", "", -1)

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: userInput,
		})

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    openai.GPT4oMini,
				Messages: messages,
			},
		)
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			return
		}

		content := resp.Choices[0].Message.Content
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: content,
		})
		fmt.Printf("O&M Expert's answer: %v\n", content)

	}
}
