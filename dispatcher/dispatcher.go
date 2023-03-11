package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"task"

	"github.com/nsqio/go-nsq"
)

type Server struct {
	ip        string
	port      string
	taskTopic string
	task      task.Task
	maxDes    int
	result    chan task.Result
}

var describes []task.TaskDes

var servers []Server

var desLock sync.Mutex

func init() {
	// 26 citys
	for i := 0; i < 8; i++ {
		describes = append(
			describes,
			task.TaskDes{
				CityId:    i,
				TypeStart: 0,
				PageStart: 0,
				PageEnd:   2,
			},
		)
	}

	servers = []Server{
		{
			ip:        "192.168.210.253",
			taskTopic: "task_queue",
			task:      task.Task{Describe: make([]task.TaskDes, 0), Goroutines: 2},
			maxDes:    4,
			result:    make(chan task.Result, 4),
		},
		{
			ip:        "192.168.210.200",
			taskTopic: "task_queue",
			task:      task.Task{Describe: make([]task.TaskDes, 0), Goroutines: 1},
			maxDes:    2,
			result:    make(chan task.Result, 2),
		},
	}
}

func resultConsumer() {

	consumer, err := nsq.NewConsumer("result_queue", "result", nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
	}
	consumer.AddHandler(nsq.HandlerFunc(resultHandle))
	if err := consumer.ConnectToNSQD("localhost:4150"); err != nil {
		log.Fatal(err)
	}
	<-consumer.StopChan
}

func resultHandle(message *nsq.Message) error {

	r := task.Result{}
	json.Unmarshal(message.Body, &r)
	fmt.Println("consumer receive result: ", r)

	for i := 0; i < len(servers); i++ {
		if servers[i].ip == r.ServerIP {
			servers[i].result <- r
		}

	}
	return nil
}

func taskDispatch() {

	var wg sync.WaitGroup

	for i := 0; len(describes) > 0; i++ {

		wg.Add(1)
		go func(server Server) {
			sendTask(server)
			wg.Done()
		}(servers[i%len(servers)])

		if (i+1)%len(servers) == 0 {
			wg.Wait()
		}
	}

	fmt.Println("finish all")
}

func sendTask(server Server) {

	desLock.Lock()
	for i := 0; i < len(describes); i++ {
		fmt.Print(describes[i].CityId, " ")
	}
	fmt.Println()
	server.task.Describe = make([]task.TaskDes, 0)
	for i := 0; i < server.maxDes; i++ {
		if len(describes) > 0 {
			server.task.Describe = append(server.task.Describe, describes[0])
			describes = describes[1:]
		}
	}
	desLock.Unlock()

	producer, err := nsq.NewProducer(server.ip+":4150", nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("send task: ", server.task)
	taskJson, _ := json.Marshal(server.task)
	if err := producer.Publish(server.taskTopic, taskJson); err != nil {
		log.Fatal("publish error: " + err.Error())
	}

	fmt.Println("waiting for result...")

	for i := 0; i < len(server.task.Describe); i++ {
		result := <-server.result
		fmt.Println("get result: ", result, i, len(server.task.Describe))
		if result.Err {
			desLock.Lock()
			des := task.TaskDes{
				CityId: result.CityId, TypeStart: result.ErrorType, PageStart: result.ErrorPage, PageEnd: result.EndPage,
			}
			describes = append(describes, des)
			desLock.Unlock()
		} else {
			fmt.Println("finish task", result, "task len", len(describes))
		}
	}
	fmt.Println("task len:", len(describes))
}

func main() {
	go taskDispatch()
	resultConsumer()

}
