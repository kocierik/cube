package main

import (
	"cube/manager"
	"cube/task"
	"cube/worker"
	"fmt"
	"log"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main() {
	// host := os.Getenv("CUBE_HOST")
	// port, _ := strconv.Atoi(os.Getenv("5555"))
	// mhost := os.Getenv("CUBE_MANAGER_HOST")
	// mport, _ := strconv.Atoi(os.Getenv("CUBE_MANAGER_PORT"))
	mHost := "localhost"
	mPort := 5556
	wHost := "localhost"
	wPort := 5555

	w := worker.Worker{
		Name:  "worker-1",
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}

	wApi := worker.Api{
		Address: wHost,
		Port:    wPort,
		Worker:  &w,
	}

	log.Printf("starting Cube worker at %s:%d", wHost, wPort)

	go w.RunTasks()
	go w.CollectStats()
	go wApi.Start()

	workers := []string{fmt.Sprintf("%s:%d", wHost, wPort)}
	m := manager.New(workers)
	mApi := manager.Api{
		Address: mHost,
		Port:    mPort,
		Manager: m,
	}

	log.Printf("starting Cube manager at %s:%d", mHost, mPort)

	go m.ProcessTasks()
	go m.UpdateTasks()

	mApi.Start()
}
