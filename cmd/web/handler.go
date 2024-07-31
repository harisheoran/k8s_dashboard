package main

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func (app *application) Roothandler(w http.ResponseWriter, request *http.Request) {
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
		app.serverError(w, err)
	}

	// K8s Client
	home, _ := os.UserHomeDir()
	configpath := filepath.Join(home, ".kube/config")
	config, err := clientcmd.BuildConfigFromFlags("", configpath)
	if err != nil {
		app.serverError(w, err)
	}
	client := kubernetes.NewForConfigOrDie(config)

	// get namespaces
	namespaceList, err := getNamespaces(client)
	if err != nil {
		app.serverError(w, err)
	}

	// get pods
	podList, err := getPods(client, selectedNamespace)
	if err != nil {
		app.serverError(w, err)
	}
	// get services
	svcList, err := getServices(client, selectedNamespace)
	if err != nil {
		app.serverError(w, err)
	}

	// get deployment
	deploymentList, err := getDeployment(client, selectedNamespace)
	if err != nil {
		app.serverError(w, err)
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
