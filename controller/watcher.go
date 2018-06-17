package controller

import (
	"fmt"
	"log"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/google/uuid"
)

// Watcher watch pod events and ensure they are running
type Watcher struct {
	kubernetes Kubernetes
	podSpec    *PodSpec
	usage      float32
	status     *Status
	minPods    int
}

// NewWatcher returns a new watcher
func NewWatcher(kubernetes Kubernetes, podSpec *PodSpec, minPods int) *Watcher {
	return &Watcher{
		kubernetes: kubernetes,
		podSpec:    podSpec,
		usage:      80,
		minPods:    minPods,
	}
}

// SetStatus sets status
func (w *Watcher) SetStatus(status *Status) {
	w.status = status
}

// Create creates namespace with numberOfPods pods
func (w *Watcher) Create(servicePort int) (pods []string, err error) {
	ns := w.podSpec.namespace

	err = w.kubernetes.CreateNamespace(ns)
	if err != nil {
		return nil, err
	}

	podNames := make([]string, w.minPods)
	for i := 0; i < w.minPods; i++ {
		name := podName(ns)
		podNames[i] = name
		err = w.kubernetes.CreatePod(ns, name, w.podSpec.image, w.podSpec.command)
		if err != nil {
			return nil, err
		}
	}

	return podNames, nil
}

func podName(namespace string) string {
	id := uuid.New().String()
	firstPart := strings.Split(id, "-")[0]
	return fmt.Sprintf("%s-%s", namespace, firstPart)
}

// Watch watch for pod events
func (w *Watcher) Watch() error {
	watcher, err := w.kubernetes.Watch(w.podSpec.namespace)
	if err != nil {
		return err
	}

	defer watcher.Stop()

	for {
		select {
		case event := <-watcher.ResultChan():
			switch obj := event.Object.(type) {
			case *v1.Pod:
				log.Printf("pod %s %s\n", event.Type, obj.GetName())
				err = w.handlePod(obj, event.Type)
				checkErr(err)
			}
		case <-w.status.Watch():
			log.Printf("pod changed status")
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
	case watch.Deleted:
		err = w.ensureNumberOfPods()
	}

	return err
}

func (w *Watcher) ensureNumberOfPods() (err error) {
	report := w.status.Report()

	log.Printf("current usage: %.2f, desired usage: %.2f", report.Usage(), w.usage)

	deltaPods := w.status.Report().Delta(w.usage)
	log.Printf("creating %d pods\n", deltaPods)

	pods, err := w.kubernetes.GetPods(w.podSpec.namespace)
	if err != nil {
		return err
	}

	if deltaPods < 0 {
		err = w.deletePods(-deltaPods, pods)
	} else if deltaPods > 0 && len(pods) < 10 {
		err = w.createPods(deltaPods)
	}

	return err
}

func (w *Watcher) createPods(numberOfPods int) error {
	for i := 0; i < numberOfPods; i++ {
		err := w.kubernetes.CreatePod(
			w.podSpec.namespace,
			podName(w.podSpec.namespace), w.podSpec.image,
			w.podSpec.command,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Watcher) deletePods(numberOfPods int, pods []v1.Pod) error {
	if len(pods)-numberOfPods < w.minPods {
		numberOfPods = len(pods) - w.minPods
	}

	for i := 0; i < numberOfPods; i++ {
		err := w.kubernetes.DeletePod(w.podSpec.namespace, pods[i].GetName())
		if err != nil {
			return err
		}
	}

	return nil
}
