package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

func main() {
	token := "sk-NUkr3DrnPJtiWZS32d667dAaCd304f20B43bAbF0D6872b21"
	config := openai.DefaultConfig(token)
	config.BaseURL = "https://vip.apiyi.com/v1"
	client := openai.NewClientWithConfig(config)

	query := "Check logs contains keyword \"Error\" and labeled with app=grafana."

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You're a Loki log querier, u could help user to analyze log, u could call multiple functions to finish tasks from user, and try to analyze the cause.",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: query,
		},
	}

	// schema
	schema := `{
        "type": "object",
        "properties": {
            "query_str": {
                "type": "string",
                "description": "Loki query string, example: {app=\"grafana\"} |= \"Error\""
            }
        },
        "required": ["query_str"]
    }`
	var schemaObj map[string]interface{}
	err := json.Unmarshal([]byte(schema), &schemaObj)
	if err != nil {
		fmt.Println("Error decoding schema:", err)
		return
	}

	tools := []openai.Tool{
		{
			Type: "function",
			Function: &openai.FunctionDefinition{
				Name:        "analyze_loki_log",
				Description: "Get logs from Loki",
				Parameters:  schemaObj,
			},
		},
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:      openai.GPT4oMini,
			Messages:   messages,
			Tools:      tools,
			ToolChoice: "auto",
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	respMsg := resp.Choices[0].Message
	toolCalls := respMsg.ToolCalls

	fmt.Printf("ChatGPT wants to call function: %v\n", toolCalls[0].Function)

}
