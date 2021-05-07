package saboter

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

type Saboter struct {
	// Kubernetes client
	Client kubernetes.Interface
	// Deletion rate / minute (only works if num victim pods > Rate to prevent cluster from falling over)
	Rate int64
}

func NewSaboter(client kubernetes.Interface, rate int64) *Saboter {
	return &Saboter{Client: client, Rate: rate}
}

func (saboter *Saboter) Start(ctx context.Context) {
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	log.Print("Registering signal handlers...")
	// signal handler
	go func() {
		log.Print("Finished registering signal handlers, listening for SIGTERM")
		<-signalChannel
		fmt.Println("Gracefully exiting")
		os.Exit(0)
	}()

	log.Print("Started saboter!")
	// Every minute find and delete saboter.Rate pods
	for range time.Tick(time.Second * 10) {
		listOptions := metav1.ListOptions{LabelSelector: "sabotage=true"}
		pods, err := saboter.Client.CoreV1().Pods("").List(context.TODO(), listOptions)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("There are %d pods in the cluster with the sabotage label\n", len(pods.Items))
	}
}
