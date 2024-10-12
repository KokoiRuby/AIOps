package main

import (
	"context"
	_ "embed"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"strings"
)

//go:embed deployment.yaml
var deployYaml string

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

	// create deployment
	// yaml → unstructured
	deployObj := &unstructured.Unstructured{}
	if err := yaml.Unmarshal([]byte(deployYaml), deployObj); err != nil {
		panic(err)
	}

	// get gvk
	apiVersion, found, err := unstructured.NestedString(deployObj.Object, "apiVersion")
	if err != nil || !found {
		log.Fatalln("apiVersion not found:", err)
	}

	kind, found, err := unstructured.NestedString(deployObj.Object, "kind")
	if err != nil || !found {
		log.Fatalln("kind not found:", err)
	}

	// to gvr
	gvr := schema.GroupVersionResource{}
	version := strings.Split(apiVersion, "/")
	if len(version) == 2 {
		gvr.Group = version[0]
		gvr.Version = version[1]
	} else {
		// res in core api group has no group, such as pod
		gvr.Version = version[0]
	}

	switch kind {
	case "Deployment":
		gvr.Resource = "deployments"
	default:
		log.Fatalf("unsupported kind: %s", kind)
	}

	_, err = dynamicClient.Resource(gvr).Namespace("default").
		Create(context.TODO(), deployObj, metav1.CreateOptions{})
	if err != nil {
		log.Fatalln(err)
	}

}

func getConfigFromSeveral() (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		nil,
	).ClientConfig()
}
