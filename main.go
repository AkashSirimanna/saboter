package main

import (
	"context"
	"flag"
	"log"
	"time"

	"akashsirimanna.com/saboter/saboter"

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

	saboter := saboter.NewSaboter(client, 1)
	saboter.Start(context.TODO())
}

func init() {
	flag.StringVar(&kubepath, "kubepath", "", "Location of kubernetes configuration, defaults to ~/.kube/config")
}
