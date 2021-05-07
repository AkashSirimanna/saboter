package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/sample-controller/pkg/signals"
)

var (
	kubepath string
)

func main() {
	flag.Parse()

	stop := signals.SetupSignalHandler()

	config, err := clientcmd.BuildConfigFromFlags("", kubepath)
	if err != nil {
		log.Fatal(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	informerFactory := kubeinformers.NewSharedInformerFactory(client, time.Second*30)
	informerFactory.Start(stop)

	listOptions := metav1.ListOptions{LabelSelector: "sabotage=true"}
	pods, err := client.CoreV1().Pods("").List(context.TODO(), listOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("There are %d pods in the cluster with the sabotage label\n", len(pods.Items))
}

func init() {
	flag.StringVar(&kubepath, "kubepath", "", "Location of kubernetes configuration, defaults to ~/.kube/config")
}
