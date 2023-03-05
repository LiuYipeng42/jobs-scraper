package main

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/nsqio/go-nsq"

	"task"
)

type Server struct {
	ip    string
	port  string
	taskTopic string
	statusTopic string
	task  task.Task
	msg   chan string
}

var describes []task.TaskDes

var servers []Server

var desLock sync.Mutex

func init() {
	// 26 citys
	for i := 0; i < 26; i++ {
		describes = append(
			describes,
			task.TaskDes{
				CityId: i,
				TypeStart: 0,
				PageStart: 0,
				PageEnd: 1,
			},
		)
	}

	servers = []Server{
		{"localhost", "4150", "task1", "status1", task.Task{Describe: make([]task.TaskDes, 2), Goroutines: 2}, make(chan string)},
		{"120.77.177.229", "4150", "task2", "status2", task.Task{Describe: make([]task.TaskDes, 2), Goroutines: 2}, make(chan string)},
	}
}

func taskStatus(server Server) {

}

func sendTask(server Server, producer nsq.Producer) {
	desLock.Lock()
	for i := 0; i < server.task.Goroutines && len(describes) > 0; i++ {
		server.task.Describe[i] = describes[0]
		describes = describes[1:]
	}
	desLock.Unlock()

	taskJson, _ := json.Marshal(server.task)
	if err := producer.Publish(server.taskTopic, taskJson); err != nil {
		log.Fatal("publish error: " + err.Error())
	}

	msg := <-server.msg

	// handle task that have error
	if msg != "success" {

	}
}

func taskDispatch() {

	producers := []nsq.Producer{}

	for i := 0; i < len(servers); i++ {
		producer, err := nsq.NewProducer(servers[i].ip+servers[i].port, nsq.NewConfig())
		if err != nil {
			log.Fatal(err)
		}
		producers = append(producers, *producer)
	}

	var wg sync.WaitGroup

	for i := 0; ; i++ {
		wg.Add(1)
		go func(server Server, producer nsq.Producer) {

			sendTask(server, producer)

			wg.Done()
		}(servers[i%len(servers)], producers[i%len(producers)])

		if (i+1)%len(servers) == 0 {
			wg.Wait()
		}
	}

}

func main() {

}
