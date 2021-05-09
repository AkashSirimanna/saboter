package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"time"

	"akashsirimanna.com/saboter/saboter"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/sample-controller/pkg/signals"
)

var (
	kubepath        string
	excludedayspath string
)

func main() {
	log.Printf("Starting saboter on %s", time.Now().Format("2006-01-02"))
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

	excludedDays := make(map[string]bool)
	if excludedayspath != "" {
		excludedDays = parseDays(excludedayspath)
	}

	saboter := saboter.NewSaboter(client, 1, excludedDays)
	saboter.Start(context.TODO())
}

func init() {
	flag.StringVar(&kubepath, "kubepath", "", "Location of kubernetes configuration, defaults to ~/.kube/config")
	flag.StringVar(&excludedayspath, "exclude", "", "Location of file containing line separated dates formatted as yyyy-mm-dd to not run saboter on")
}

func parseDays(path string) map[string]bool {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	excludedDays := make(map[string]bool)

	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		dateStr := fileScanner.Text()
		if _, err := time.Parse("2006-01-02", dateStr); err != nil {
			log.Fatal(err)
		}
		excludedDays[dateStr] = true
	}

	return excludedDays
}
