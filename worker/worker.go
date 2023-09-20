package worker

import (
	"cube/task"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/client"
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

func (w *Worker) RunTask() {
	fmt.Println("Start and stop the task")
}

func (w *Worker) StartTask() {
	fmt.Println("Start the task")
}

func (w *Worker) StopTask(t task.Task) task.DockerResult {
	c := task.Config{
		Name:  "test-container-1",
		Image: "postgres:13",
		Env: []string{
			"POSTGRES_USER=cube",
			"POSTGRES_PASSWORD=secret",
		},
	}
	dc, _ := client.NewClientWithOpts(client.FromEnv)
	d := task.Docker{
		Client: dc,
		Config: c,
	}
	result := d.Stop(t.ID.String())

	if result.Error != nil {
		log.Printf("Error stopping container %v: %v", d.ContainerId, result.Error.Error())
	}
	t.FinishTime = time.Now().UTC()
	t.State = task.Completed
	w.Db[t.ID] = t
	log.Printf("Stopped and removed container %v for task %v", d.ContainerId, t.Name)
	return result

}
