package taskschedule

type TaskStatus int

type Task interface {
	ID() int
	Next() bool
	Status() TaskStatus
	Callback() func(Task) error
}
