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

type Pod struct {
	Name      string
	State     string
	CreatedAt time.Time
}

type Namspace struct {
	Name string
}

type Svc struct {
	Name string
	Type string
}

type Deployment struct {
	Name    string
	Replica *int32
}

type K8sDashData struct {
	Pods              []Pod
	Svcs              []Svc
	Deployment        []Deployment
	Namespace         []Namspace
	SelectedNamespace string
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
	selectedNamespace := "default"
	payloadNamespace := request.URL.Query().Get("namespace")
	if payloadNamespace != "" {
		selectedNamespace = payloadNamespace
	}

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

	// get namespaces
	namespaceList, err := getNamespaces(client)
	if err != nil {
		fmt.Println(err.Error())
	}

	// get pods
	podList, err := getPods(client, selectedNamespace)
	if err != nil {
		fmt.Println(err.Error())
	}
	// get services
	svcList, err := getServices(client, selectedNamespace)
	if err != nil {
		fmt.Println(err.Error())
	}

	// get deployment
	deploymentList, err := getDeployment(client, selectedNamespace)
	if err != nil {
		fmt.Println(err.Error())
	}

	data := K8sDashData{
		Pods:              podList,
		Svcs:              svcList,
		Deployment:        deploymentList,
		Namespace:         namespaceList,
		SelectedNamespace: selectedNamespace,
	}

	template.Execute(w, data)
}

// get pods
func getPods(client *kubernetes.Clientset, selectedNamespace string) ([]Pod, error) {
	podList := []Pod{}
	pods, err := client.CoreV1().Pods(selectedNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return podList, err
	}

	for _, thispod := range pods.Items {
		pod := Pod{
			Name:      thispod.Name,
			State:     string(thispod.Status.Phase),
			CreatedAt: thispod.CreationTimestamp.Time,
		}
		podList = append(podList, pod)
	}

	return podList, nil
}

// get services
func getServices(client *kubernetes.Clientset, selectedNamespace string) ([]Svc, error) {
	servicesList := []Svc{}
	services, err := client.CoreV1().Services(selectedNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return servicesList, err
	}

	for _, thissvc := range services.Items {
		svc := Svc{
			Name: thissvc.Name,
			Type: string(thissvc.Spec.Type),
		}
		servicesList = append(servicesList, svc)
	}

	return servicesList, nil
}

// get deployemet
func getDeployment(client *kubernetes.Clientset, selectedNamespace string) ([]Deployment, error) {
	deploymentList := []Deployment{}

	deployments, err := client.AppsV1().Deployments(selectedNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return deploymentList, err
	}

	for _, thisdeployment := range deployments.Items {
		deployment := Deployment{
			Name:    thisdeployment.Name,
			Replica: thisdeployment.Spec.Replicas,
		}
		deploymentList = append(deploymentList, deployment)
	}
	return deploymentList, nil
}

// get namespaces
func getNamespaces(client *kubernetes.Clientset) ([]Namspace, error) {
	namespacesList := []Namspace{}

	namespaces, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return namespacesList, err
	}

	for _, thisNamespace := range namespaces.Items {
		namespace := Namspace{
			Name: thisNamespace.Name,
		}
		namespacesList = append(namespacesList, namespace)
	}

	return namespacesList, nil
}
