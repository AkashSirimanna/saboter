package saboter

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

type Saboter struct {
	// Kubernetes client
	Client kubernetes.Interface
	// Deletion rate / minute (only works if num victim pods > Rate to prevent cluster from falling over)
	Interval int64
	// Deletion rate / minute (only works if num victim pods > Rate to prevent cluster from falling over)
	Rate int64
	//Splice of days to not run saboter
	ExcludedDays map[string]bool
}

func NewSaboter(client kubernetes.Interface, interval, rate int64, excludedDays map[string]bool) *Saboter {
	return &Saboter{Client: client, Interval: interval, Rate: rate, ExcludedDays: excludedDays}
}

func (saboter *Saboter) Start(ctx context.Context) {
	gracePeriodSeconds := int64(1)
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
	for range time.Tick(time.Minute * time.Duration(saboter.Interval)) {
		if today := time.Now().Format("2006-01-02"); saboter.ExcludedDays[time.Now().Format("2006-01-02")] == true {
			log.Printf("Skipping sabotages on %s", today)
			time.Sleep(24*time.Hour + 1*time.Second) //Suspend execution for just over 1 day and continue
			continue
		}
		listOptions := metav1.ListOptions{LabelSelector: "sabotage=true"}
		pods, err := saboter.Client.CoreV1().Pods("").List(ctx, listOptions)
		if err != nil {
			log.Fatal(err)
		}

		if len(pods.Items) == 0 {
			log.Println("No pods with sabotage label, waiting until pods with the label appear")
			labelExists := make(chan bool, 1)
			go func() {
				for range time.Tick(time.Second * 10) {
					pods, err := saboter.Client.CoreV1().Pods("").List(ctx, listOptions)
					if err != nil {
						log.Fatal(err)
					}

					if len(pods.Items) != 0 {
						labelExists <- true
						break
					}
				}
			}()
			select {
			case <-labelExists:
				continue
			case <-time.After(10 * time.Minute):
				log.Fatal("No pods with sabotage label appeared in 10 minutes, goodbye")
			}
		}

		numCandidatesToSabotage := math.Min(float64(len(pods.Items)), float64(saboter.Rate))

		for i := 0; i < int(numCandidatesToSabotage); i++ {
			randomIndex := rand.Int() % len(pods.Items)
			candidate := pods.Items[randomIndex]
			pods.Items[randomIndex] = pods.Items[len(pods.Items)-1]
			pods.Items[len(pods.Items)-1] = v1.Pod{}
			pods.Items = pods.Items[:len(pods.Items)-1]

			log.Printf("Sabotaging pod %s on namespace %s running on node %s...", candidate.Name, candidate.Namespace, candidate.Spec.NodeName)
			if saboter.Client.CoreV1().Pods(candidate.Namespace).Delete(ctx, candidate.Name, metav1.DeleteOptions{GracePeriodSeconds: &gracePeriodSeconds}) != nil {
				log.Println(err)
				continue
			}
			log.Printf("Finished sabotaging pod %s on namespace %s running on node %s", candidate.Name, candidate.Namespace, candidate.Spec.NodeName)
		}
	}
}
