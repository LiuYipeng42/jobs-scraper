package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/nsqio/go-nsq"
	"log"
	"pojo"
)

var db *sqlx.DB
var err error

var selectCompanySql string
var selectJobSql string
var updateJobSql string
var insertCompanySql string
var insertJobSql string

func init() {
	db, err = sqlx.Open("mysql", "root:121522734a@tcp(localhost:3306)/jobs")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}

	selectCompanySql = "select id as cid from companys where name = ?"

	var buffer bytes.Buffer
	buffer.WriteString("select j.id, j.name, position, c.name as cname ")
	buffer.WriteString("from jobs j left join companys c on j.company_id = c.id ")
	buffer.WriteString("where j.name = ? and j.position = ? and cname = ?")
	selectJobSql = buffer.String()

	updateJobSql = "update jobs set salary=? experience=? degree=? tags=? describe=? where id = ?"

	insertCompanySql = "insert into companys (name, type, size, main_business, `describe`) values (?, ?, ?, ?, ?)"

	insertJobSql = "insert into jobs (name, salary, position, experience, degree, tags, `describe`, url, company_id) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
}

func startConsumer(channel string) {
	cfg := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("jobs", channel, cfg)
	if err != nil {
		log.Fatal(err)
	}
	// 设置消息处理函数
	consumer.AddHandler(nsq.HandlerFunc(process))
	// 连接到单例nsqd
	if err := consumer.ConnectToNSQD("localhost:4150"); err != nil {
		log.Fatal(err)
	}
	<-consumer.StopChan
}

func process(message *nsq.Message) error {
	job := pojo.Job{}
	json.Unmarshal(message.Body, &job)
	// if strings.Contains(job.Url, "51job") {
	// 	log.Println(job)
	// }
	// if strings.Contains(job.Url, "58.com") {
	// 	log.Println(job)
	// }
	var jobs []pojo.Job
	if db.Select(&jobs, selectJobSql, job.Name, job.Position, job.CName); err != nil {
		return err
	}
	if len(jobs) > 0 {
		for i := 0; i < len(jobs); i++ {
			if job.Url == jobs[i].Url {
				r, err := db.Exec(updateJobSql, job.Salary, job.Experience, job.Experience, job.Tags, job.Describe, jobs[i].Id)
				if err != nil {
					fmt.Println(r, err)
				}
			}
		}
	} else {
		r, err := db.Exec(insertCompanySql, job.CName, job.CType, job.CSize, job.MainBusiness, job.CDescribe)
		var id int64
		if err == nil {
			id, _ = r.LastInsertId()
			fmt.Println("insert company:", id, job.CName)
		} else {
			company := pojo.Company{}
			db.Select(&company, selectCompanySql, job.CName)
		}

		r, err = db.Exec(insertJobSql, job.CName, job.Salary, job.Position, job.Experience, job.Degree, job.Tags, job.Describe, job.Url, id)
		id, _ = r.LastInsertId()
		fmt.Println("insert job: ", id, job.Name, job.CName)
		if err != nil {
			fmt.Println(r, err)
		}
	}

	return nil
}

func main() {

	startConsumer("job_channel")

}
