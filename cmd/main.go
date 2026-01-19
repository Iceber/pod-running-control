package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	celtypes "github.com/google/cel-go/common/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	plugincel "k8s.io/apiserver/pkg/admission/plugin/cel"
	"k8s.io/apiserver/pkg/cel/environment"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/Iceber/pod-running-control/cel"
)

type runningGate struct {
	gvr       schema.GroupVersionResource
	namespace string
	name      string
}

func main() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	/*
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
	*/

	var kubeconfig string
	var gateGVR string
	var gate runningGate
	var expression string
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig file with authorization and master location information.")

	flag.StringVar(&gateGVR, "gate-gvr", "", "group and version of gate resource")
	flag.StringVar(&gate.namespace, "gate-namespace", "", "namespace of gate resource")
	flag.StringVar(&gate.name, "gate-name", "", "name of gate resource")
	flag.StringVar(&expression, "gate-expression", "", "name of gate resource")
	flag.Parse()

	if gvr, _ := schema.ParseResourceArg(gateGVR); gvr != nil {
		gate.gvr = *gvr
	}
	if gate.gvr.Empty() {
		gate.gvr = schema.GroupVersionResource{Group: "pod-running-control.io", Version: "v1alpha1", Resource: "podrunninggates"}
		expression = "size(object.spec.gates) == 0"
	}

	compositionEnvTemplate, err := plugincel.NewCompositionEnv(plugincel.VariablesTypeName, environment.MustBaseEnvSet(environment.DefaultCompatibilityVersion()))
	if err != nil {
		panic(err)
	}
	optionalVars := plugincel.OptionalVariableDeclarations{HasParams: false, HasAuthorizer: false}
	compiler := plugincel.NewCompiler(compositionEnvTemplate.EnvSet)
	condition := compiler.CompileCELExpression(&cel.ValidationCondition{Expression: expression}, optionalVars, environment.StoredExpressions)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		os.Exit(1)
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		os.Exit(1)
	}

	informer := dynamicinformer.NewFilteredDynamicInformer(client, gate.gvr, gate.namespace, 0, nil, func(opts *metav1.ListOptions) {
		opts.FieldSelector = "metadata.name=" + gate.name
	}).Informer()

	gateChecker := func(obj *unstructured.Unstructured) {
		result, err := cel.Evaluate(context.TODO(), obj, condition)
		if err != nil {
			fmt.Println(err)
			return
		}
		if result.Error != nil {
			fmt.Println(result.Error)
			return
		}
		if result.EvalResult == celtypes.True {
			os.Exit(0)
		}
	}
	if _, err = informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { gateChecker(obj.(*unstructured.Unstructured)) },
		UpdateFunc: func(_, obj interface{}) { gateChecker(obj.(*unstructured.Unstructured)) },
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
