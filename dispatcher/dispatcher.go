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
	task  task.Task
	result   chan task.Result
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
		{"localhost", "4150", "task_queue1", task.Task{Describe: make([]task.TaskDes, 2), Goroutines: 2}, make(chan task.Result)},
		{"120.77.177.229", "4150", "task_queue2", task.Task{Describe: make([]task.TaskDes, 2), Goroutines: 2}, make(chan task.Result)},
	}
}

// receive result from server
func resultConsumer() {

	consumer, err := nsq.NewConsumer("result_queue", "result", nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
	}
	// 设置消息处理函数
	consumer.AddHandler(nsq.HandlerFunc(resultHandle))
	// 连接到单例nsqd
	if err := consumer.ConnectToNSQD("localhost:4150"); err != nil {
		log.Fatal(err)
	}
	<-consumer.StopChan

}

func resultHandle(message *nsq.Message) error {
	
	r := task.Result{}
	json.Unmarshal(message.Body, &r)

	for i := 0; i < len(servers); i++ {
		if servers[i].ip == r.ServerIP {
			servers[i].result <- r
		}
	}

	return nil
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

	result := <-server.result

	// handle task that have error
	if result.Err {
		desLock.Lock()
		des := task.TaskDes{
			CityId: result.CityId, TypeStart: result.ErrorType, PageStart: result.ErrorPage, PageEnd: result.EndPage,
		}
		describes = append(describes, des)
		desLock.Unlock()
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
