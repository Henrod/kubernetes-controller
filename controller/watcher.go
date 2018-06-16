package controller

import (
	"fmt"
	"log"
	"strings"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/google/uuid"
)

// Watcher watch pod events and ensure they are running
type Watcher struct {
	kubernetes   Kubernetes
	namespace    string
	image        string
	ticker       time.Duration
	numberOfPods int
}

// NewWatcher returns a new watcher
func NewWatcher(
	kubernetes Kubernetes,
	namespace, image string,
	numberOfPods int,
	ticker time.Duration,
) *Watcher {
	return &Watcher{kubernetes, namespace, image, ticker, numberOfPods}
}

// Create creates namespace with numberOfPods pods
func (w *Watcher) Create(servicePort int) error {
	ns := w.namespace

	err := w.kubernetes.CreateNamespace(ns)
	if err != nil {
		return err
	}

	for i := 0; i < w.numberOfPods; i++ {
		err = w.kubernetes.CreatePod(ns, podName(ns), w.image)
		if err != nil {
			return err
		}
	}

	err = w.kubernetes.CreateService(ns, ns, servicePort)
	if err != nil {
		return err
	}

	return nil
}

func podName(namespace string) string {
	id := uuid.New().String()
	firstPart := strings.Split(id, "-")[0]
	return fmt.Sprintf("%s-%s", namespace, firstPart)
}

// Watch watch for pod events
func (w *Watcher) Watch() error {
	watcher, err := w.kubernetes.Watch(w.namespace)
	if err != nil {
		return err
	}

	defer watcher.Stop()

	ticker := time.NewTicker(w.ticker)

	for {
		select {
		case event := <-watcher.ResultChan():
			switch obj := event.Object.(type) {
			case *v1.Pod:
				log.Printf("pod %s %s\n", event.Type, obj.GetName())
				err = w.handlePod(obj, event.Type)
				checkErr(err)
			}
		case <-ticker.C:
			err = w.ensureNumberOfPods()
			checkErr(err)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		log.Printf("error: %q\n", err)
	}
}

func (w *Watcher) handlePod(pod *v1.Pod, eventType watch.EventType) (err error) {
	switch eventType {
	case watch.Modified:
		err = w.ensureNumberOfPods()
	}

	return err
}

func (w *Watcher) ensureNumberOfPods() error {
	pods, err := w.kubernetes.GetPods(w.namespace)
	if err != nil {
		return err
	}

	deltaPods := len(pods) - w.numberOfPods
	if deltaPods < 0 {
		err = w.createPods(-deltaPods)
	} else if deltaPods > 0 {
		err = w.deletePods(deltaPods, pods)
	}

	return err
}

func (w *Watcher) createPods(numberOfPods int) error {
	for i := 0; i < numberOfPods; i++ {
		err := w.kubernetes.CreatePod(w.namespace, podName(w.namespace), w.image)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Watcher) deletePods(numberOfPods int, pods []v1.Pod) error {
	for i := 0; i < numberOfPods; i++ {
		err := w.kubernetes.DeletePod(w.namespace, pods[i].GetName())
		if err != nil {
			return err
		}
	}

	return nil
}
