package controller

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// Kubernetes has methods to access kubernetes
type Kubernetes interface {
	CreatePod(namespace, name, image string) error
	CreateService(namespace, name string, port int) error
	CreateNamespace(namespace string) error

	GetPods(namespace string) ([]v1.Pod, error)

	DeletePod(namespace, name string) error

	Watch(namespace string) (watch.Interface, error)
}
