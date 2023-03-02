package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"pojo"
	"strconv"
	"strings"
	"sync"
	"time"
	"utils"

	"github.com/nsqio/go-nsq"
	"github.com/tebeka/selenium"
)

func switchCity(chrome utils.ChromeSerADri, city int) {
	t := time.Now()
	allCity := chrome.WaitAndFindOne("div.allcity", 5, 1)
	allCity.Click()

	table := chrome.WaitAndFindOne("div.el-dialog", 5, 1)

	citys, _ := table.FindElements(selenium.ByCSSSelector, "div.grid-item>span")
	citys[city].Click()

	yes, _ := table.FindElement(selenium.ByCSSSelector, "div.el-dialog__footer>span")
	yes.Click()

	fmt.Println(time.Since(t))
}

func getInfo(html string) (job pojo.Job) {

	name := utils.RegExpFindOne(html, "title=.*?class=\"jname at\"")
	job.Name = name[7:len(name)-18]

	// update := utils.RegExpFindOne(html, "class=\"time\">.*?</span>")
	// job.Update = update[13:len(update)-7]

	jobInfo := strings.Split(utils.RegExpFindOne(html, "class=\"info\">.*?</p>"), "<span")
	job.Salary = jobInfo[1][32:len(jobInfo[1])-7]
	job.Position = jobInfo[3][20:len(jobInfo[3])-7]
	job.Position = strings.ReplaceAll(job.Position, "·", "-")
	if len(jobInfo) > 5 {
		job.Experience = jobInfo[5][20:len(jobInfo[5])-7]
		if len(jobInfo) > 7 {
			job.Degree = jobInfo[7][20:len(jobInfo[7])-18]
		}
	}

	tagElems := strings.Split(utils.RegExpFindOne(html, "class=\"tags\">.*?</p>"), "title=")[1:]

	var tagBuffer bytes.Buffer
	for _, tag := range tagElems {
		tagBuffer.WriteString(tag[1:strings.Index(tag, ">")-1] + " ")
	}
	job.Tags = tagBuffer.String()
	job.Tags = job.Tags[:len(job.Tags)]

	url := utils.RegExpFindOne(html, "<a .*? href=\".*?\" target=\"_blank\" class=\"el\">")
	job.Url = url[strings.Index(url, "http"):len(url) - 29]

	cname := utils.RegExpFindOne(html, "class=\"cname at\">.*?</a>")
	job.CName = cname[17:len(cname) - 4]

	companyInfo := utils.RegExpFindOne(html, "class=\"dc at\">.*?</p>")
	CInfo := strings.Split(companyInfo[14:len(companyInfo) - 4], " | ")
	job.Company.CType = CInfo[0]
	if len(CInfo) > 1 {
		job.Company.CSize = CInfo[1]
	}

	business := utils.RegExpFindOne(html, "class=\"int at\">.*?</p>")
	job.MainBusiness = strings.ReplaceAll(business[15:len(business) - 4], "/", " ")

	return
}

func saveJobInfo(chrome utils.ChromeSerADri, firstCity, lastCity int) {

	chrome.Webdriver.Get("https://we.51job.com/pc/search")

	city := firstCity
	switchCity(chrome, city)
	
	pageNum := 0
	for {
		jobs := chrome.WaitAndFindAll("div.j_joblist>div[sensorsname]", 5)
		fmt.Println(len(jobs))

		for _, job := range jobs {
			fmt.Printf("%d %d ", city, pageNum)
			html, _ := job.GetAttribute("outerHTML")
			job := getInfo(html)
			fmt.Println(job.Name)

			jobJson, _ := json.Marshal(job)
			if err := producer.Publish("jobs", jobJson); err != nil {
				log.Fatal("publish error: " + err.Error())
			}

		}

		pageNumS, _ := chrome.WaitAndFindOne("li.number.active", 2, 1).Text()
		pageNum, _ = strconv.Atoi(pageNumS)
		if pageNum < 200 {
			chrome.WaitAndFindOne("button.btn-next", 5, 1).Click()
		} else {
			city += 1
			if city > lastCity {
				return
			}
			switchCity(chrome, city)
		}
	}
}

var producer *nsq.Producer

func init() {
	cfg := nsq.NewConfig()
	var err error
	producer, err = nsq.NewProducer("localhost:4150", cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	chrome := utils.InitClientByDriver("./chromedriver", 8080, false)
	wd := chrome.Webdriver
	defer chrome.Service.Stop()
	defer wd.Quit()

	var wg sync.WaitGroup

	// for i := 0; i < 1; i++ {
	// 	wg.Add(1)
	// 	go func(first, last int) {
	// 		chrome := utils.InitClientByRemote("http://localhost:4444/wd/hub")
	// 		wd := chrome.Webdriver
	// 		defer wd.Quit()

	// 		saveJobInfo(chrome, first, last)
	// 		wg.Done()
	// 	}(80 * i, 80 * i + 80)
	// }

	wg.Add(1)
	go func() {
		// chrome := utils.InitClientByRemote("http://localhost:4444/wd/hub")
		// wd := chrome.Webdriver
		// defer wd.Quit()

		saveJobInfo(chrome, 400, 430)
		wg.Done()
	}()

	wg.Wait()
}
