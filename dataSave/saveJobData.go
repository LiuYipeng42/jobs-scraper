package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"job"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/nsqio/go-nsq"
)

var db *sqlx.DB
var err error

var selectJob string
var updateJob string
var insertJob string
var selectCompany string
var updateCompany string
var insertCompany string

func init() {
	db, err = sqlx.Open("mysql", "root:121522734a@tcp(localhost:3306)/jobs")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}

	selectCompany = "select id as cid from companys where name = ?"

	var buffer bytes.Buffer
	buffer.WriteString("select id, company_id as cid from jobs where ")
	buffer.WriteString("name=? and type=? and salary=? and position=? and experience=? and degree=?")
	selectJob = buffer.String()

	updateJob = "update jobs set tags=?, `describe`=? where id=?"

	updateCompany = "update companys set type=?, size=?, main_business=?, describe where id=?"

	insertCompany = "insert into companys (name, type, size, main_business, `describe`) values (?, ?, ?, ?, ?)"

	insertJob = "insert into jobs (name, type, salary, position, experience, degree, tags, `describe`, url, company_id) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
}

func startConsumer(channel string) {
	consumer, err := nsq.NewConsumer("jobs", channel, nsq.NewConfig())
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
	j := job.Job{}
	json.Unmarshal(message.Body, &j)

	// fmt.Println(j.Name)
	var jobs []job.Job
	if err = db.Select(&jobs, selectJob, j.Name, j.Type, j.Salary, j.Position, j.Experience, j.Degree); err != nil {
		return err
	}
	var company []job.Company
	if err = db.Select(&company, selectCompany, j.CName); err != nil {
		return err
	}

	// exist company
	if len(company) > 0 {
		fmt.Println("update company: ", j.CName)
		db.Exec(updateCompany, company[0].CType, company[0].CSize, company[0].MainBusiness, company[0].CId)

		if len(jobs) > 0 {
			// exist job
			for i := 0; i < len(jobs); i++ {
				if jobs[i].CId == company[0].CId {
					fmt.Println("update job: ", jobs[i].Id, j.Name, j.CName)
					if j.Url == jobs[i].Url {
						_, err := db.Exec(updateJob, j.Tags, j.Describe, jobs[i].Id)
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		} else {
			// new job
			fmt.Println("insert job: ", j.Name, j.CName)
			db.Exec(insertJob, j.Name, j.Type, j.Salary, j.Position, j.Experience, j.Degree, j.Tags, j.Describe, j.Url, company[0].CId)
		}
	}

	// new job and company
	if len(company) == 0 {
		fmt.Println("insert company: ", j.CName)
		r, _ := db.Exec(insertCompany, j.CName, j.CType, j.CSize, j.MainBusiness, j.CDescribe)
		cid, _ := r.LastInsertId()
		fmt.Println("insert job: ", j.Name, j.CName)
		db.Exec(insertJob, j.Name, j.Type, j.Salary, j.Position, j.Experience, j.Degree, j.Tags, j.Describe, j.Url, cid)
	}

	return nil
}

func main() {

	startConsumer("job_channel")

}
