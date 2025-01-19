package taskschedule

type TaskStatus int

type Task interface {
	Next() bool
	Status() TaskStatus
}
