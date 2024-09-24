package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
)

func main() {
	ctx := context.Background()

	// llm
	llm, err := openai.New(
		openai.WithModel("gpt-4o-mini"),
		openai.WithBaseURL("https://vip.apiyi.com/v1"),
		openai.WithToken("BLOCKED"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// mem
	mem := memory.ConversationBuffer{}

	// chain
	chain := chains.NewConversation(llm, &mem)

	// predict
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("-> ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.Replace(userInput, "\n", "", -1)

		m := make(map[string]any)
		m["input"] = userInput

		// how to pass input?
		resp, err := chains.Predict(ctx, chain, m)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Chain response is: %v\n", resp)
		fmt.Printf("Current context is: %v\n", chain.Memory)

	}
}
