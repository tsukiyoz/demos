package taskschedule

type TaskStatus int

const (
	StatusUnknown TaskStatus = iota
	StatusPending
	StatusDownloadCompleted
	StatusAnalysisCompleted
	StatusAllCompleted
)

func (s TaskStatus) Next() (TaskStatus, bool) {
	switch s {
	case StatusPending:
		return StatusDownloadCompleted, true
	case StatusDownloadCompleted:
		return StatusAnalysisCompleted, true
	case StatusAnalysisCompleted:
		return StatusAllCompleted, true
	case StatusAllCompleted:
		return StatusAllCompleted, false
	default:
		return StatusUnknown, false
	}
}

type Task struct {
	ID       int
	Status   TaskStatus
	Payload  map[string]string
	Callback func(*Task) error
}
