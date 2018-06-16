package clientset

import (
	"fmt"

	homedir "github.com/mitchellh/go-homedir"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//Clientset connects to minikube and returns the clientset
func Clientset() (kubernetes.Interface, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	kubeConfig := fmt.Sprintf("%s/.kube/config", home)

	var config *rest.Config
	config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
