package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type pod struct {
	Name      string
	State     string
	CreatedAt time.Time
}

type svc struct {
	Name string
	Type string
}

type deployment struct {
	Name    string
	Replica *int32
}

type k8sDashData struct {
	Pods       []pod
	Svcs       []svc
	Deployment []deployment
}

func main() {
	// create a router
	router := http.NewServeMux()

	// router map the path with the handler
	router.HandleFunc("/", roothandler)

	// start the web server
	fmt.Println("Starting the web server at port 1313")
	err := http.ListenAndServe(":1313", router)
	if err != nil {
		log.Fatal(err)
	}
}

func roothandler(w http.ResponseWriter, request *http.Request) {
	myTemplates := []string{
		"./ui/home.page.gohtml",
	}

	template, err := template.ParseFiles(myTemplates...)
	if err != nil {
		log.Println("Unable to parse html templates", err)
	}

	// K8s Client
	home, _ := os.UserHomeDir()
	configpath := filepath.Join(home, ".kube/config")
	config, err := clientcmd.BuildConfigFromFlags("", configpath)
	if err != nil {
		fmt.Println("ERROR: Unable to authenticate, check your kubeconfig files ", err.Error())
	}
	client := kubernetes.NewForConfigOrDie(config)

	// get pods
	podList, err := getPods(client)
	if err != nil {
		fmt.Println(err.Error())
	}
	// get services
	svcList, err := getServices(client)
	if err != nil {
		fmt.Println(err.Error())
	}

	// get deployment
	deploymentList, err := getDeployment(client)

	data := k8sDashData{
		Pods:       podList,
		Svcs:       svcList,
		Deployment: deploymentList,
	}

	template.Execute(w, data)
}

// get pods
func getPods(client *kubernetes.Clientset) ([]pod, error) {
	podList := []pod{}
	Namespace := "default"
	pods, err := client.CoreV1().Pods(Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return podList, err
	}

	for _, thispod := range pods.Items {
		pod := pod{
			Name:      thispod.Name,
			State:     string(thispod.Status.Phase),
			CreatedAt: thispod.CreationTimestamp.Time,
		}
		podList = append(podList, pod)
	}

	return podList, nil
}

// get services
func getServices(client *kubernetes.Clientset) ([]svc, error) {
	Namespace := "default"
	servicesList := []svc{}
	services, err := client.CoreV1().Services(Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return servicesList, err
	}

	for _, thissvc := range services.Items {
		svc := svc{
			Name: thissvc.Name,
			Type: string(thissvc.Spec.Type),
		}
		servicesList = append(servicesList, svc)
	}

	return servicesList, nil
}

func getDeployment(client *kubernetes.Clientset) ([]deployment, error) {
	namspace := "default"
	deploymentList := []deployment{}

	deployments, err := client.AppsV1().Deployments(namspace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return deploymentList, err
	}

	for _, thisdeployment := range deployments.Items {
		deployment := deployment{
			Name:    thisdeployment.Name,
			Replica: thisdeployment.Spec.Replicas,
		}
		deploymentList = append(deploymentList, deployment)
	}
	return deploymentList, nil
}
