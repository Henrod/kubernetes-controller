package controller

import "log"

// Status has the pod statuses
const (
	StatusReady      = "ready"
	StatusOccupied   = "occupied"
	StatusTerminated = "terminated"
)

// Status holds information about pod statuses
type Status struct {
	podsStatuses map[string]string
	channel      chan struct{}
}

// NewStatus creates a Status with statuses ready
func NewStatus(pods []string) *Status {
	status := &Status{
		podsStatuses: map[string]string{},
		channel:      make(chan struct{}),
	}

	for _, pod := range pods {
		status.podsStatuses[pod] = StatusReady
	}

	return status
}

// Update updates pod statuses
func (s *Status) Update(name, status string) {
	log.Printf("updating status on %s to %s\n", name, status)

	defer func() {
		s.channel <- struct{}{}
	}()

	if status == StatusTerminated {
		delete(s.podsStatuses, name)
		return
	}

	s.podsStatuses[name] = status
}

// Report returns how many pods per status
func (s *Status) Report() *Report {
	report := &Report{}

	for _, status := range s.podsStatuses {
		switch status {
		case StatusReady:
			report.ReadyPods = report.ReadyPods + 1
		case StatusOccupied:
			report.OccupiedPods = report.OccupiedPods + 1
		}
	}

	return report
}

// Watch returns a channel that receives on everyy update
func (s *Status) Watch() chan struct{} {
	return s.channel
}
