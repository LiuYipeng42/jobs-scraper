package task


type Task struct {
	Describe []TaskDes
	Goroutines int
}

type TaskDes struct {
	CityId int
	TypeStart int
	PageStart int
	PageEnd int
}
