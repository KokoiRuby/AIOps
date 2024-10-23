/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/KokoiRuby/k8scopilot/cmd/utils"
	"github.com/sashabaranov/go-openai"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spf13/cobra"
)

// eventCmd represents the event command
var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		eventLog, err := getPodEventsAndLgs()
		if err != nil {
			fmt.Println("Error getting pod events log")
		}
		fmt.Println(eventLog)

		// chatgpt
		res, err := sendToChatGPT(eventLog)
		if err != nil {
			fmt.Println("Error sending to chat")
		}
		fmt.Println(res)

	},
}

func sendToChatGPT(podEventLog map[string][]string) (string, error) {
	client, err := utils.NewOpenAI()
	if err != nil {
		return "", err
	}

	combined := "Please find warning event & log from below pod: \n"
	for podName, logs := range podEventLog {
		combined += fmt.Sprintf("\tPod: %s\n", podName)
		for _, log := range logs {
			combined += fmt.Sprintf("%s\n", log)
		}
		combined += "\n"
	}

	fmt.Print(combined)

	msg := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You're a K8s expert, please help diagnose what kind of issues user is facing",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Below is the aggregation of warning type logs of pod(s): \n%s\nPlease provide constructive suggestions & feasible commands", combined),
		},
	}

	resp, err := client.Client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini,
		Messages: msg,
	})
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil

}

// string: []string → pod: logs
func getPodEventsAndLgs() (map[string][]string, error) {
	clientGo, err := utils.NewClientGo(kubeconfig)
	if err != nil {
		return nil, err
	}
	res := make(map[string][]string)

	// Warning
	events, err := clientGo.ClientSet.CoreV1().Events(namespace).List(context.Background(), metav1.ListOptions{
		FieldSelector: "type=Warning",
	})
	if err != nil {
		return nil, err
	}

	for _, event := range events.Items {
		podName := event.InvolvedObject.Name
		namespace := event.InvolvedObject.Namespace
		msg := event.Message

		if event.InvolvedObject.Kind == "Pod" {
			logOption := &corev1.PodLogOptions{}
			req := clientGo.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, logOption)
			podLogs, err := req.Stream(context.Background())
			if err != nil {
				return nil, err
			}
			defer func(podLogs io.ReadCloser) {
				err := podLogs.Close()
				if err != nil {

				}
			}(podLogs)

			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(podLogs)
			if err != nil {
				// if no log entry in pod
				continue
			}
			// if had log entry
			res[podName] = append(res[podName], fmt.Sprintf("Event: %s", msg))
			fmt.Println(msg)
			res[podName] = append(res[podName], fmt.Sprintf("Namespace: %s", namespace))
			fmt.Println(namespace)
			res[podName] = append(res[podName], fmt.Sprintf("Logs: %s", buf.String()))
			fmt.Println(buf.String())
		}
	}
	return res, nil
}

func init() {
	analyzeCmd.AddCommand(eventCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// eventCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// eventCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
