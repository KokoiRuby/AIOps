package main

import (
	"context"
	_ "embed"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Dynamic Client
	config, err := getConfigFromSeveral()
	if err != nil {
		panic(err.Error())
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	gvr := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}

	unStructObj, err := dynamicClient.Resource(gvr).Namespace("kube-system").
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	podList := &corev1.PodList{}
	// unstructured to podList
	if err = runtime.DefaultUnstructuredConverter.
		FromUnstructured(unStructObj.UnstructuredContent(), podList); err != nil {
		panic(err)
	}

	for _, pod := range podList.Items {
		fmt.Println(pod.Name)
	}

}

func getConfigFromSeveral() (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		nil,
	).ClientConfig()
}
