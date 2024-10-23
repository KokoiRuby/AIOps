/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/KokoiRuby/k8scopilot/cmd/utils"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/restmapper"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// chatgptCmd represents the chatgpt command
var chatgptCmd = &cobra.Command{
	Use:   "chatgpt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		startChat()
	},
}

func init() {
	askCmd.AddCommand(chatgptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatgptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatgptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// 1. input as query to llm
func startChat() {
	// from stdin
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("I'm K8s Copilot, anything I could do for u?")

	for {
		fmt.Print("> ")
		if scanner.Scan() {
			input := scanner.Text()
			if input == "exit" {
				fmt.Println("Bye!")
				break
			}
			if input == "" {
				continue
			}
			//fmt.Println("Your query is: ", text)
			fmt.Println(processInput(input))
		}
	}
}

var systemPrompt = `
You are a Kubernetes Copilot. Your task is to assist users in completing their Kubernetes-related tasks.
Whenever your response includes YAML, please output it in YAML format only, and don't put it into YAML (code) block.'
`

// 2. process input by sending chat msg to llm given systemPrompt & input
func processInput(input string) string {
	client, err := utils.NewOpenAI()
	if err != nil {
		return err.Error()
	}
	//resp, err := client.SendMessage(systemPrompt, input)
	//if err != nil {
	//	return err.Error()
	//}
	resp := funcCalling(input, client)
	return resp
}

// 3. define function to call
func funcCalling(input string, client *utils.OpenAI) string {
	// gen YAML & deploy
	f1 := openai.FunctionDefinition{
		Name:        "generateAndDeployResource",
		Description: "Generate K8s YAML manifest & deploy to cluster",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"user_input": {
					Type:        jsonschema.String,
					Description: "Extract verb, resource and necessary flags",
				},
			},
			Required: []string{"user_input"},
		},
	}
	t1 := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f1,
	}

	// get K8s res
	f2 := openai.FunctionDefinition{
		Name:        "getResource",
		Description: "Get K8s resources",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"namespace": {
					Type:        jsonschema.String,
					Description: "Namespace where resource is",
				},
				"resource": {
					Type:        jsonschema.String,
					Description: "K8s built-in resource, for example: pods, deployments, services, you can also use singular or short name (if had)",
				},
			},
			Required: []string{"namespace", "resource"},
		},
	}
	t2 := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f2,
	}

	// delete K8s res
	f3 := openai.FunctionDefinition{
		Name:        "deleteResource",
		Description: "Delete K8s resources",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"namespace": {
					Type:        jsonschema.String,
					Description: "Namespace where resource is",
				},
				"resource": {
					Type:        jsonschema.String,
					Description: "K8s built-in resource, for example: pods, deployments, services, you can also use singular or short name (if had)",
				},
				"resource_name": {
					Type:        jsonschema.String,
					Description: "Name of the resource to be deleted",
				},
			},
			Required: []string{"namespace", "resource", "resource_name"},
		},
	}
	t3 := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f3,
	}

	dialogue := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		},
	}

	resp, err := client.Client.CreateChatCompletion(context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4oMini,
			Messages: dialogue,
			Tools:    []openai.Tool{t1, t2, t3},
		},
	)
	if err != nil {
		return err.Error()
	}

	msg := resp.Choices[0].Message
	if len(msg.ToolCalls) != 1 {
		return fmt.Sprintf("No appropriate tool is found, %v", len(msg.ToolCalls))
	}

	// build chat history
	dialogue = append(dialogue, msg)
	//return fmt.Sprintf("Function to call: %s, arg: %s", msg.ToolCalls[0].Function.Name, msg.ToolCalls[0].Function.Arguments)
	result, err := invokeFunc(client, msg.ToolCalls[0].Function.Name, msg.ToolCalls[0].Function.Arguments)
	if err != nil {
		return err.Error()
	}
	return result
}

func generateAndDeployResource(client *utils.OpenAI, input string) (string, error) {
	sysPromt := `
You're a K8s resource generator.
Please generate K8s YAML based on user input.
Don't include it into YAML code block.
`
	// gen YAML
	yamlContent, err := client.SendMessage(sysPromt, input)
	if err != nil {
		return "", err
	}
	//return yamlContent, nil

	// client-go
	clientGo, err := utils.NewClientGo(kubeconfig)
	if err != nil {
		return "", err
	}
	res, err := restmapper.GetAPIGroupResources(clientGo.DiscoveryClient)
	if err != nil {
		return "", err
	}
	// YAML to Unstructured
	unstructuredObj := &unstructured.Unstructured{}
	_, _, err = scheme.Codecs.UniversalDeserializer().Decode([]byte(yamlContent), nil, unstructuredObj)
	if err != nil {
		return "", err
	}

	// mapper from gvr to gvk
	mapper := restmapper.NewDiscoveryRESTMapper(res)
	gvk := unstructuredObj.GroupVersionKind()
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return "", err
	}

	namespace := unstructuredObj.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}

	// create
	_, err = clientGo.DynamicClient.Resource(mapping.Resource).Namespace(namespace).Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Resource %s created successfully", unstructuredObj.GetName()), nil
}

func getResource(namespace, resource string) (string, error) {
	clientGo, err := utils.NewClientGo(kubeconfig)
	if err != nil {
		return "", err
	}
	resource = strings.ToLower(resource)

	var gvr schema.GroupVersionResource
	switch resource {
	case "pods", "pod", "po":
		gvr = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	case "services", "service", "svc":
		gvr = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	case "deployments", "deployment":
		gvr = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	default:
		return "", fmt.Errorf("resource %s not supported", resource)
	}

	resList, err := clientGo.DynamicClient.Resource(gvr).Namespace(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	result := ""
	for _, item := range resList.Items {
		result += fmt.Sprintf("%s\n", item.GetName())
	}

	return result, nil

}

func deleteResource() (string, error) {
	return "del res", nil
}

// 4. invoke func
func invokeFunc(client *utils.OpenAI, name, args string) (string, error) {
	switch name {
	case "generateAndDeployResource":
		params := struct {
			UserInput string `json:"user_input"`
		}{}
		if err := json.Unmarshal([]byte(args), &params); err != nil {
			return "", err
		}
		return generateAndDeployResource(client, params.UserInput)
	case "getResource":
		params := struct {
			Namespace string `json:"namespace"`
			Resource  string `json:"resource"`
		}{}
		if err := json.Unmarshal([]byte(args), &params); err != nil {
			return "", err
		}
		return getResource(params.Namespace, params.Resource)
	case "deleteResource":
		return deleteResource()
	default:
		return "", fmt.Errorf("unknown function %s", name)
	}
}
