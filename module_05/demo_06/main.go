package main

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"time"
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

	// shared informer fac
	informerFactory := informers.NewSharedInformerFactory(clientset, time.Hour*12)

	// build informer & watch
	deployInformer := informerFactory.Apps().V1().Deployments()
	informer := deployInformer.Informer()
	// Lister contains local cache
	deployLister := deployInformer.Lister()
	_, err = informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAddDeployment,
		UpdateFunc: onUpdateDeployment,
		DeleteFunc: onDeleteDeployment,
	})
	if err != nil {
		return
	}

	// signal to stop
	stopper := make(chan struct{})
	defer close(stopper)
	// start informer & wait for cache to sync
	// start will send a signal to chan
	informerFactory.Start(stopper)
	informerFactory.WaitForCacheSync(stopper)

	deployments, err := deployLister.Deployments("default").List(labels.Everything())
	if err != nil {
		panic(err.Error())
	}
	for _, deploy := range deployments {
		fmt.Println(deploy.Name)
	}

	// receive from chan, main goroutine gets blocked until chan is closed.
	<-stopper
}

func onAddDeployment(obj interface{}) {
	deploy := obj.(*appsv1.Deployment)
	fmt.Println("add a deployment:", deploy.Name)
}

func onUpdateDeployment(old, new interface{}) {
	oldDeploy := old.(*appsv1.Deployment)
	newDeploy := new.(*appsv1.Deployment)
	fmt.Println("update deployment:", oldDeploy.Name, newDeploy.Name)
}

func onDeleteDeployment(obj interface{}) {
	deploy := obj.(*appsv1.Deployment)
	fmt.Println("delete a deployment:", deploy.Name)
}

func getConfigFromSeveral() (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		nil,
	).ClientConfig()
}
