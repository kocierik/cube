package manager

import (
	"cube/task"
	"fmt"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Manager struct {
	Pending       queue.Queue
	TaskDb        map[string][]task.Task
	EventDb       map[string][]task.TaskEvent
	Workers       []string
	WorkerTaskMap map[string][]uuid.UUID
	TaskWorkerMap map[uuid.UUID]string
}

func (m *Manager) SelectWorker() {
	fmt.Println("Select an appropriate worker")
}
func (m *Manager) UpdateTasks() {
	fmt.Println("Update tasks")
}
func (m *Manager) SendWork() {
	fmt.Println("Send work to workers")
}
