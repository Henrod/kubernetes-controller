package controller

// Report has the number of ready and occupied pods
type Report struct {
	ReadyPods    int
	OccupiedPods int
}

// Usage returns the usage report
func (r *Report) Usage() float32 {
	ready := float32(r.ReadyPods)
	occupied := float32(r.OccupiedPods)
	total := ready + occupied

	return occupied / total
}

// Delta returns how many pods to create or to delete to
// reach expected usage
func (r *Report) Delta(expectedUsage float32) int {
	/*
		(occupied + x) / (ready + occupied + x) = expectedUsage = eu
		occupied + x = eu*ready + eu*occupied + eu*x
		x*(1-eu) = eu*ready + (eu - 1)*occupied
		x*(1 - eu) = eu*ready - (1 - eu)*occupied
		x = eu*ready/(1 - eu) - occupied
	*/

	ready := float32(r.ReadyPods)
	occupied := float32(r.OccupiedPods)

	delta := expectedUsage*ready/(1-expectedUsage) - occupied
	return int(delta)
}
