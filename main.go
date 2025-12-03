package main

import (
	"custom-controller/controller"
	"flag"
	"fmt"
	"log"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/Users/ashokmaurya/.kube/config", "Please provide the path to the KUBECONFIG file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %s\n", err.Error())
		fmt.Println("Switching to Alternative way")
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("Error building kubeconfig: %s\n", err.Error())
		}
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error building clientset: %s\n", err.Error())
	}
	ch := make(chan struct{})
	sharedInformer := informers.NewSharedInformerFactory(clientSet, 10*time.Minute)
	ctrl := controller.NewController(clientSet, sharedInformer.Apps().V1().Deployments())
	sharedInformer.Start(ch)
	ctrl.Start(ch)
}
