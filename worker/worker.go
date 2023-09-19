package worker

import (
	"cube/task"
	"fmt"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Worker struct {
	Name      string
	Queue     queue.Queue
	Db        map[uuid.UUID]task.Task
	TaskCount int
}

func (w *Worker) CollectStats() {
	fmt.Println("Collect stats")
}

func (w *Worker) runTask() {
	fmt.Println("Start and stop the task")
}

func (w *Worker) startTask() {
	fmt.Println("Start the task")
}

func (w *Worker) stopTask() {
	fmt.Println("Stop the task")
}
