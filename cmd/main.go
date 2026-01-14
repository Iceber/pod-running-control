package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/Iceber/pod-running-control/api/v1alpha1"
	clientset "github.com/Iceber/pod-running-control/client-go/clientset/versioned"
	informers "github.com/Iceber/pod-running-control/client-go/informers/externalversions/api/v1alpha1"
)

func main() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	namespace := os.Getenv("POD_RUNNING_GATE_NAMESPACE")
	if namespace == "" {
		fmt.Fprintln(os.Stderr, "POD_RUNNING_GATE_NAMESPACE is empty, please set it and restart.")
		<-signalCh
		os.Exit(1)
	}

	gateName := os.Getenv("POD_RUNNING_GATE_NAME")
	if gateName == "" {
		fmt.Fprintln(os.Stderr, "POD_RUNNING_GATE_NAME is empty, please set it and restart.")
		<-signalCh
		os.Exit(1)
	}

	var kubeconfig string
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig file with authorization and master location information.")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		os.Exit(1)
	}
	client, err := clientset.NewForConfig(config)
	if err != nil {
		os.Exit(1)
	}

	informer := informers.NewFilteredPodRunningGateInformer(client, namespace, 0, nil, func(opts *metav1.ListOptions) {
		opts.FieldSelector = "metadata.name=" + gateName
	})

	gateChecker := func(obj *v1alpha1.PodRunningGate) {
		if len(obj.Spec.Gates) != 0 {
			return
		}
		os.Exit(0)
	}
	if _, err = informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { gateChecker(obj.(*v1alpha1.PodRunningGate)) },
		UpdateFunc: func(_, obj interface{}) { gateChecker(obj.(*v1alpha1.PodRunningGate)) },
		DeleteFunc: func(_ interface{}) { /*log*/ },
	}); err != nil {
		os.Exit(1)
	}

	stopCh := make(chan struct{})
	go func() {
		<-signalCh
		close(stopCh)
	}()

	informer.Run(stopCh)
}
