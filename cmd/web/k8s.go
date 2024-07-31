package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

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
