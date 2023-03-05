package task


type Task struct {
	Jobs []Job
	Goroutines int
}

type Job struct {
	TypeId int
	CityId int
	PageStart int
	PageEnd int
}
