package controller

// PodSpec defines a pod
type PodSpec struct {
	namespace string
	image     string
	command   []string
}

// NewPodSpec returns a pod spec
func NewPodSpec(namespace, image string, command []string) *PodSpec {
	return &PodSpec{namespace, image, command}
}
