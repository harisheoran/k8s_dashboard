package main

import "time"

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
