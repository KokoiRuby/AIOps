package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
)

func main() {
	// csv → jsonl
	// convertCSVToJSONList("./data/log.csv", "./log.jsonl")

	ctx := context.Background()

	token := "sk-NUkr3DrnPJtiWZS32d667dAaCd304f20B43bAbF0D6872b21"
	config := openai.DefaultConfig(token)
	config.BaseURL = "https://vip.apiyi.com/v1"
	client := openai.NewClientWithConfig(config)

	// upload jsonl to llm with purpose
	file, err := client.CreateFile(ctx, openai.FileRequest{
		FilePath: "./log.jsonl",
		Purpose:  "fine-tune",
	})
	if err != nil {
		fmt.Printf("Upload JSONL file error: %v\n", err)
		return
	}

	// create tuning job by file id & model given
	fineTuningJob, err := client.CreateFineTuningJob(ctx, openai.FineTuningJobRequest{
		TrainingFile: file.ID,
		Model:        "gpt-4o-2024-08-06", // https://platform.openai.com/docs/guides/fine-tuning
	})
	if err != nil {
		fmt.Printf("Creating new fine tune model error: %v\n", err)
		return
	}
	fmt.Printf("Fine tuning job id is: %v\n", fineTuningJob.ID)

	// it shall take a while

	// retrieve job by job id
	fineTuningJob, err = client.RetrieveFineTuningJob(ctx, fineTuningJob.ID)
	if err != nil {
		fmt.Printf("Getting fine tune model error: %v\n", err)
		return
	}
	fmt.Println(fineTuningJob.Status)
	fmt.Println(fineTuningJob.FineTunedModel)

	model_id := "ft:gpt-4o-mini-2024-07-18:he3::9tbKC81t"

	// use tuned model
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model_id,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a log level identifier, please identify the log level for each log entry, output P0, P1 or P2 directly.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Disk I/O error",
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

func convertCSVToJSONList(csvPath, jsonLPath string) {
	// input stream
	csvFile, err := os.Open(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)
	// col name → idx
	// returned row in csv is []string
	header, err := reader.Read()
	if err != nil {
		log.Fatal(err)
	}
	colIdx := make(map[string]int)
	for i, colName := range header {
		colIdx[colName] = i
	}

	// output stream
	jsonLFile, err := os.Create(jsonLPath)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonLFile.Close()
	encoder := json.NewEncoder(jsonLFile)
	encoder.SetIndent("", "  ")

	for {
		row, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatal(err)
		}

		logPrint := row[colIdx["log"]]
		priority := row[colIdx["priority"]]
		jsonLEntry := []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a log level identifier, please identify the log level f each log entry, output P0, P1 or P2 directly.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: logPrint,
			},
			{
				Role:    openai.ChatMessageRoleAssistant,
				Content: priority,
			},
		}

		if err = encoder.Encode(jsonLEntry); err != nil {
			log.Fatal(err)
		}
	}
}
