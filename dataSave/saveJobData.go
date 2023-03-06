package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"job"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/nsqio/go-nsq"
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

	fmt.Println(j)

	// var jobs []job.Job
	// if db.Select(&jobs, selectJobSql, j.Name, j.Position, j.CName); err != nil {
	// 	return err
	// }

	// if len(jobs) > 0 {
	// 	for i := 0; i < len(jobs); i++ {
	// 		fmt.Println("exist job: ", jobs[i].Id, j.Name, j.CName)
	// 		if j.Url == jobs[i].Url {
	// 			r, err := db.Exec(updateJobSql, j.Salary, j.Experience, j.Experience, j.Tags, j.Describe, jobs[i].Id)
	// 			if err != nil {
	// 				fmt.Println(r, err)
	// 			}
	// 		}
	// 	}
	// } else {
	// 	var company []job.Company
	// 	db.Select(&company, selectCompanySql, j.CName)
	// 	var cid int64
	// 	if len(company) == 0 {
	// 		r, _ := db.Exec(insertCompanySql, j.CName, j.CType, j.CSize, j.MainBusiness, j.CDescribe)
	// 		if r != nil {
	// 			cid, _ = r.LastInsertId()
	// 			fmt.Println("insert company:", cid, j.CName)
	// 		}
	// 	} else {
	// 		cid = int64(company[0].CId)
	// 		fmt.Println("exist company:", cid, j.CName)
	// 	}

	// 	var jid int64
	// 	r, _ := db.Exec(insertJobSql, j.CName, j.Salary, j.Position, j.Experience, j.Degree, j.Tags, j.Describe, j.Url, cid)
	// 	if r != nil {
	// 		jid, _ = r.LastInsertId()
	// 	} else {
	// 		fmt.Println(err)
	// 	}
	// 	fmt.Println("insert job: ", jid, j.Name, j.CName)

	// }

	return nil
}

func main() {

	startConsumer("job_channel")

}
