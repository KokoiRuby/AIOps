package main

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := getConfigFromSeveral()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	timeOut := int64(30)
	watcher, _ := clientset.CoreV1().Pods("default").
		Watch(context.Background(), metav1.ListOptions{TimeoutSeconds: &timeOut})

	// iter watch chan
	for event := range watcher.ResultChan() {
		item := event.Object.(*corev1.Pod)

		switch event.Type {
		case watch.Added:
			fmt.Printf("pod added: %s\n", item.Name)
		case watch.Modified:
			fmt.Printf("pod modified: %s\n", item.Name)
		case watch.Deleted:
			fmt.Printf("pod deleted: %s\n", item.Name)
		case watch.Error:
			fmt.Printf("pod error: %s\n", item.Name)
		default:
			fmt.Printf("event: %s\n", event.Type)
		}
	}
}

func getConfigFromSeveral() (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		nil,
	).ClientConfig()
}
