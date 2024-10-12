package main

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type Controller struct {
	indexer  cache.Indexer
	queue    workqueue.TypedRateLimitingInterface[string]
	informer cache.Controller
}

// NewController constructor
func NewController(queue workqueue.TypedRateLimitingInterface[string], indexer cache.Indexer, informer cache.Controller) *Controller {
	return &Controller{
		informer: informer,
		indexer:  indexer,
		queue:    queue,
	}
}

func (c *Controller) processNextItem() bool {
	key, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Done(key)

	err := c.syncToStdout(key)
	c.handleErr(err, key)
	return true
}

func (c *Controller) syncToStdout(key string) error {
	// get obj from indexer by key
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		fmt.Printf("Fetching object with key %s from store failed with %v\n", key, err)
		return err
	}

	if !exists {
		fmt.Printf("Deployment %s does not exist anymore\n", key)
	} else {
		deployment := obj.(*appsv1.Deployment)
		fmt.Printf("Sync/Add/Update for Deployment %s, Replicas: %d\n", deployment.Name, *deployment.Spec.Replicas)
		// simulate err
		if deployment.Name == "test-deployment" {
			time.Sleep(2 * time.Second)
			return fmt.Errorf("simulated error for deployment %s", deployment.Name)
		}
	}
	return nil
}

func (c *Controller) handleErr(err error, key string) {
	if err == nil {
		c.queue.Forget(key)
		return
	}

	if c.queue.NumRequeues(key) < 5 {
		fmt.Printf("Retry %d for key %s\n", c.queue.NumRequeues(key), key)
		// add it back to queue
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	fmt.Printf("Dropping deployment %q out of the queue: %v\n", key, err)
}

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

	// rate limit queue
	queue := workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]())

	// build informer & watch
	deployInformer := informerFactory.Apps().V1().Deployments()
	informer := deployInformer.Informer()
	_, err = informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { onAddDeployment(obj, queue) },
		UpdateFunc: func(old, new interface{}) { onUpdateDeployment(new, queue) },
		DeleteFunc: func(obj interface{}) { onDeleteDeployment(obj, queue) },
	})
	if err != nil {
		return
	}

	controller := NewController(queue, deployInformer.Informer().GetIndexer(), informer)

	// signal to stop
	stopper := make(chan struct{})
	defer close(stopper)
	// start informer & wait for cache to sync
	// start will send a signal to chan
	informerFactory.Start(stopper)
	informerFactory.WaitForCacheSync(stopper)

	// process item in controller
	go func() {
		for {
			if !controller.processNextItem() {
				break
			}
		}
	}()

	// receive from chan, main goroutine gets blocked until chan is closed.
	<-stopper
}

func onAddDeployment(obj interface{}, queue workqueue.TypedRateLimitingInterface[string]) {
	// get key from obj
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err == nil {
		queue.Add(key)
	}
}

func onUpdateDeployment(new interface{}, queue workqueue.TypedRateLimitingInterface[string]) {
	key, err := cache.MetaNamespaceKeyFunc(new)
	if err == nil {
		queue.Add(key)
	}
}

func onDeleteDeployment(obj interface{}, queue workqueue.TypedRateLimitingInterface[string]) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err == nil {
		queue.Add(key)
	}
}

func getConfigFromSeveral() (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		nil,
	).ClientConfig()
}
