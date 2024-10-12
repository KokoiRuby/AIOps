package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func main() {
	// parse kubectl
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s get <resource>\n", os.Args[0])
		os.Exit(1)
	}
	command := os.Args[1]
	kind := os.Args[2]

	if command != "get" {
		fmt.Println("Unsupported command:", command)
		os.Exit(1)
	}

	// Dynamic Client & ClientSet for discovery
	config, err := getConfigFromSeveral()
	if err != nil {
		panic(err.Error())
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	discoveryClient := clientset.Discovery()
	// get api resources
	apiGroupResources, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		panic(err)
	}
	// build mapper from resources
	mapper := restmapper.NewDiscoveryRESTMapper(apiGroupResources)

	gvk := schema.GroupVersionKind{
		Group:   "mygroup.example.com",
		Version: "v1alpha1",
		Kind:    kind,
	}

	// ++mapping
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		panic(err)
	}

	resourceInterface := dynamicClient.Resource(mapping.Resource).Namespace("default")
	resources, err := resourceInterface.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, resource := range resources.Items {
		fmt.Println(resource.GroupVersionKind().String())
	}

}

func getConfigFromSeveral() (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		nil,
	).ClientConfig()
}
