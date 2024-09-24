package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
)

func main() {
	token := "BLOCKED"
	baseURL := "https://vip.apiyi.com/v1"
	client := buildClient(token, baseURL)

	command1 := "Help me modify config of service \"gateway\", config name is: vendor, modifiy its value to alipay."
	command2 := "Help me restart service \"gateway\"."
	command3 := "Help me deploy a Kubernetes Deployment, where image is nginx."
	messages1 := getMessages(command1)
	messages2 := getMessages(command2)
	messages3 := getMessages(command3)
	messagesList := [][]openai.ChatCompletionMessage{messages1, messages2, messages3}

	schema1, err := buildSchemaFromPath("./schema1.json")
	if err != nil {
		log.Fatal(err)
	}
	schema2, err := buildSchemaFromPath("./schema2.json")
	if err != nil {
		log.Fatal(err)
	}
	schema3, err := buildSchemaFromPath("./schema3.json")
	if err != nil {
		log.Fatal(err)
	}

	tools := []openai.Tool{
		{
			Type: "function",
			Function: &openai.FunctionDefinition{
				Name:        "modifyConf",
				Description: "Modify configuration.",
				Parameters:  schema1,
			},
		},
		{
			Type: "function",
			Function: &openai.FunctionDefinition{
				Name:        "restartService",
				Description: "Restart Service.",
				Parameters:  schema2,
			},
		},
		{
			Type: "function",
			Function: &openai.FunctionDefinition{
				Name:        "applyManifest",
				Description: "Apply manifest.",
				Parameters:  schema3,
			},
		},
	}

	for _, msg := range messagesList {
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:      openai.GPT4oMini,
				Messages:   msg,
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

}

func buildSchemaFromPath(path string) (map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	jsonByte, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var schema map[string]interface{}
	err = json.Unmarshal(jsonByte, &schema)
	return schema, err
}

func getMessages(cmd string) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You're a command executor, u could help user to run comand, u could call multiple functions to finish tasks from user, and try to analyze the cause.",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: cmd,
		},
	}
}

func buildClient(token, baseURL string) *openai.Client {
	config := openai.DefaultConfig(token)
	config.BaseURL = baseURL
	return openai.NewClientWithConfig(config)
}
