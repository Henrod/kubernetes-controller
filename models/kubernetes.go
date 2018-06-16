package models

import (
	"strings"

	"github.com/Henrod/kube-controller/clientset"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Kubernetes access Kubernetes and implements controller.Kubernetes
type Kubernetes struct {
	client kubernetes.Interface
}

// NewKubernetes connects to kubernetes and returns
func NewKubernetes() (*Kubernetes, error) {
	client, err := clientset.Clientset()
	if err != nil {
		return nil, err
	}

	return &Kubernetes{client}, nil
}

// CreatePod creates pod
func (k *Kubernetes) CreatePod(namespace, name, image string) error {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app":  namespace,
				"name": name,
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{Name: "main", Image: image},
			},
		},
	}

	_, err := k.client.CoreV1().Pods(namespace).Create(pod)

	return err
}

// CreateNamespace creates namesapce on kubernetes
func (k *Kubernetes) CreateNamespace(namespace string) error {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	_, err := k.client.CoreV1().Namespaces().Create(ns)
	if err != nil && !isAlreadyExistsError(err) {
		return err
	}

	return nil
}

// CreateService creates a service of type NodePort on kubernetes
func (k *Kubernetes) CreateService(namespace, name string, port int) error {
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{Port: int32(port)},
			},
			Type: v1.ServiceTypeNodePort,
			Selector: map[string]string{
				"app": namespace,
			},
		},
	}

	_, err := k.client.CoreV1().Services(namespace).Create(service)
	if err != nil && !isAlreadyExistsError(err) {
		return err
	}

	return nil
}

// Watch watch events on namespace
func (k *Kubernetes) Watch(namespace string) (watch.Interface, error) {
	opts := metav1.ListOptions{}
	watcher, err := k.client.CoreV1().Pods(namespace).Watch(opts)
	return watcher, err
}

// GetPods returns the pods on namespace
func (k *Kubernetes) GetPods(namespace string) ([]v1.Pod, error) {
	opts := metav1.ListOptions{
		LabelSelector: labels.Set(map[string]string{
			"app": namespace,
		}).String(),
	}

	pods, err := k.client.CoreV1().Pods(namespace).List(opts)
	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}

// DeletePod deletes pod with name within namespace
func (k *Kubernetes) DeletePod(namespace, name string) error {
	opts := &metav1.DeleteOptions{}
	err := k.client.CoreV1().Pods(namespace).Delete(name, opts)
	return err
}

func isAlreadyExistsError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "already exists")
}
