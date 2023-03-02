package main

import (
	"encoding/json"
	"fmt"
	"os"
	"pojo"
	"strings"
	"testing"
	"utils"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/nsqio/go-nsq"
)

func TestGoquery(t *testing.T) {
	content, err := os.ReadFile("../mainPage.html")
	if err != nil {
		panic(err)
	}
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(content)))

	dom.Find("div.j_joblist>div[sensorsname]>a").Each(func(i int, element *goquery.Selection) {
		fmt.Println(element.Attr("href"))
	})
}

// go test --run=TestRe
func TestRe(t *testing.T) {
	content, err := os.ReadFile("../mainPage.html")
	if err != nil {
		panic(err)
	}
	// (?:.|\n)

	html := string(content)

	data := utils.RegExpFindOne(html, "<a .*? href=\".*?\" target=\"_blank\" class=\"el\">")
	fmt.Println(data[strings.Index(data, "https") : len(data)-29])

}

func TestNsq(t *testing.T) {

	cfg := nsq.NewConfig()
	producer, _ := nsq.NewProducer("localhost:4150", cfg)

	job := pojo.Job{Name: "1", Update: "2", Salary: "3", Position: "4", Experience: "5", Degree: "6", Tags: "7",
		Describe: "8", Url: "9", Company: pojo.Company{CName: "10", CType: "11", CSize: "12", MainBusiness: "13", CDescribe: "14"}}

	jobJson, _ := json.Marshal(job)
	producer.Publish("jobs", jobJson)

}

func TestMysql(t *testing.T) {
	db, err := sqlx.Open("mysql", "root:121522734a@tcp(localhost:3306)/jobs")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}

	// sql := "insert into companys (name) values (?)"
	// r, err := db.Exec(sql, "江西一铁科技有限公司（江西一铁（美团合作商））")
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(r.LastInsertId())
	// }

	// sql := "select "
	// sql += "j.id, j.name, j.update, salary, position, experience, degree, tags, j.describe as jdescribe, url, "
	// sql += "c.name as cname, type, size, main_business, c.describe as cdescribe "
	// sql += "from jobs j left join companygo test --run=Tests c on j.company_id = c.id "
	// sql += "where j.name = ?"
	// var jobs []pojo.Job
	// err = db.Select(&jobs, sql, "送餐员")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(jobs)

	// sql := "update jobs set name= ? where id = ?"
	// r, err := db.Exec(sql, 123, 1)
	// if err != nil {
	// 	fmt.Println(r, err)
	// }

	var company pojo.Company
	db.Select(&company, "select id as cid from company where name = ?", "上海金蝶网络科技有限公司")

	fmt.Println(company.CId)
}
