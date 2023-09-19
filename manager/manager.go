package manager

import (
	"cube/task"
	"fmt"
	"github.com/google/uuid"
)

type Manager struct {
	Pending       queue.Queue
	TaskDb        map[string][]Task
	EventDb       map[string][]TaskEvent
	Worker        []string
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
