package main

import (
	"cube/manager"
	"cube/task"
	"cube/worker"
	"fmt"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main() {
	// host := os.Getenv("CUBE_HOST")
	// port, _ := strconv.Atoi(os.Getenv("5555"))
	// mhost := os.Getenv("CUBE_MANAGER_HOST")
	// mport, _ := strconv.Atoi(os.Getenv("CUBE_MANAGER_PORT"))
	mhost := "localhost"
	mport := 5556
	whost := "localhost"
	wport := 5555

	fmt.Println("Starting Cube worker")

	w1 := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	wapi1 := worker.Api{Address: whost, Port: wport, Worker: &w1}
	w2 := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	wapi2 := worker.Api{Address: whost, Port: wport + 1, Worker: &w2}
	w3 := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	wapi3 := worker.Api{Address: whost, Port: wport + 2, Worker: &w3}

	go w1.RunTasks()
	go w1.UpdateTasks()
	go wapi1.Start()

	go w2.RunTasks()
	go w2.UpdateTasks()
	go wapi2.Start()

	go w3.RunTasks()
	go w3.UpdateTasks()
	go wapi3.Start()

	fmt.Println("Starting Cube manager")

	workers := []string{
		fmt.Sprintf("%s:%d", whost, wport),
		fmt.Sprintf("%s:%d", whost, wport+1),
		fmt.Sprintf("%s:%d", whost, wport+2),
	}
	m := manager.New(workers, "epvm")
	mapi := manager.Api{Address: mhost, Port: mport, Manager: m}

	go m.ProcessTasks()
	go m.UpdateTasks()
	go m.DoHealthChecks()

	mapi.Start()
}
