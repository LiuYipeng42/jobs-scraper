package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"job"
	"log"
	"strings"
	"sync"
	"task"
	"time"
	"utils"

	"github.com/nsqio/go-nsq"
	"github.com/tebeka/selenium"
)

func switchCity(chrome utils.ChromeSerADri, city int) {

	allCity := chrome.WaitAndFindOne("div.allcity", 10, 1)
	allCity.Click()

	fliter := chrome.WaitAndFindOne("div.j_filter", 5, 1)
	tables, _ := fliter.FindElements(selenium.ByCSSSelector, "div[role=dialog]")
	table := tables[0]

	tagTbale, _ := table.FindElement(selenium.ByCSSSelector, "div.tags-text")
	tags, _ := tagTbale.FindElements(selenium.ByCSSSelector, "i")

	for i := 0; i < len(tags); i++ {
		tags[i].Click()
	}

	// time.Sleep(2000 * time.Millisecond)
	citys, _ := table.FindElements(selenium.ByCSSSelector, "div.grid-item>span")
	citys[city].Click()

	yes, _ := table.FindElement(selenium.ByCSSSelector, "div.el-dialog__footer>span")
	yes.Click()

}

func switchJobType(chrome utils.ChromeSerADri, jobId int) {

	selectType := chrome.WaitAndFindOne("div.e_e.e_com", 10, 1)
	selectType.Click()

	fliter := chrome.WaitAndFindOne("div.j_filter", 5, 1)
	tables, _ := fliter.FindElements(selenium.ByCSSSelector, "div[role=dialog]")
	table := tables[1]

	// 清除已有的 tag
	tagTbale, _ := table.FindElement(selenium.ByCSSSelector, "div.tags-text")
	time.Sleep(1 * time.Second)
	tags, _ := tagTbale.FindElements(selenium.ByCSSSelector, "span>i")
	for i := 0; i < len(tags); i++ {
		tags[i].Click()
	}

	// 找到 jobId 的对应的 type
	tabList, _ := table.FindElements(selenium.ByCSSSelector, "div[role=tablist]>div")
	panelList, _ := table.FindElements(selenium.ByCSSSelector, "div[role=tabpanel]")

	tabLen := []int{15, 19, 24, 37, 42, 49, 52, 56, 59, 62, 70, 75}

	var tab selenium.WebElement
	var panel selenium.WebElement
	relativeId := jobId

	for i := 0; i < len(tabLen); i++ {
		if jobId-tabLen[i] < 0 {
			tab = tabList[i]
			panel = panelList[i]
			if i > 0 {
				relativeId = jobId - tabLen[i-1]
			}
			break
		}
	}

	tab.Click()
	jobTypes, _ := panel.FindElements(selenium.ByCSSSelector, "div.table-body-tr-td>span")
	jobTypes[relativeId].Click()
	all, _ := panel.FindElement(selenium.ByCSSSelector, "div.clickAll>span")
	all.Click()
	yes, _ := table.FindElement(selenium.ByCSSSelector, "div.el-dialog__footer>span")
	yes.Click()
	search, _ := chrome.Webdriver.FindElement(selenium.ByCSSSelector, "button#search_btn")
	search.Click()

}

func parser(html string) (j job.Job) {

	name := utils.RegExpFindOne(html, "title=.*?class=\"jname at\"")
	j.Name = name[7 : len(name)-18]

	jobInfo := strings.Split(utils.RegExpFindOne(html, "class=\"info\">.*?</p>"), "<span")
	j.Salary = jobInfo[1][32 : len(jobInfo[1])-7]
	j.Position = jobInfo[3][20 : len(jobInfo[3])-7]
	j.Position = strings.ReplaceAll(j.Position, "·", "-")
	if len(jobInfo) > 5 {
		j.Experience = jobInfo[5][20 : len(jobInfo[5])-7]
		if len(jobInfo) > 7 {
			j.Degree = jobInfo[7][20 : len(jobInfo[7])-18]
		}
	}

	tagElems := strings.Split(utils.RegExpFindOne(html, "class=\"tags\">.*?</p>"), "title=")[1:]

	var tagBuffer bytes.Buffer
	for _, tag := range tagElems {
		tagBuffer.WriteString(tag[1:strings.Index(tag, ">")-1] + " ")
	}
	j.Tags = tagBuffer.String()
	j.Tags = j.Tags[:len(j.Tags)]

	url := utils.RegExpFindOne(html, "<a .*? href=\".*?\" target=\"_blank\" class=\"el\">")
	j.Url = url[strings.Index(url, "http") : len(url)-29]

	cname := utils.RegExpFindOne(html, "class=\"cname at\">.*?</a>")
	j.CName = cname[17 : len(cname)-4]

	companyInfo := utils.RegExpFindOne(html, "class=\"dc at\">.*?</p>")
	CInfo := strings.Split(companyInfo[14:len(companyInfo)-4], " | ")
	j.Company.CType = CInfo[0]
	if len(CInfo) > 1 {
		j.Company.CSize = CInfo[1]
	}

	business := utils.RegExpFindOne(html, "class=\"int at\">.*?</p>")
	j.MainBusiness = strings.ReplaceAll(business[15:len(business)-4], "/", " ")

	return
}

func sendJobData(chrome utils.ChromeSerADri, jobProducer *nsq.Producer, resultProducer *nsq.Producer, des task.TaskDes) {

	jobType := des.TypeStart
	page := des.PageStart
	// send error msg
	defer sendResult(resultProducer, des, jobType, page)

	chrome.Webdriver.Get("https://we.51job.com/pc/search")

	switchCity(chrome, des.CityId)
	
	for ; jobType < 75; jobType++ {
		switchJobType(chrome, jobType)
		for ; page < des.PageEnd; page++ {
			jobs := chrome.WaitAndFindAll("div.j_joblist>div[sensorsname]", 5)
			fmt.Printf("city: %d, type: %d page: %d jobs: %d\n", des.CityId, jobType, page, len(jobs))
	
			for _, job := range jobs {
				html, _ := job.GetAttribute("outerHTML")
				job := parser(html)
				fmt.Println(job.Name)
	
				jobJson, _ := json.Marshal(job)
				if err := jobProducer.Publish("jobs", jobJson); err != nil {
					log.Fatal("publish error: " + err.Error())
				}
			}
	
			input := chrome.WaitAndFindOne("input#jump_page", 2, 1)
			input.Clear()
			input.SendKeys(fmt.Sprintf("%d", page+1))
			button, _ := chrome.Webdriver.FindElement(selenium.ByCSSSelector, "span.jumpPage")
			button.Click()
		}
	}
}

func sendResult(producer *nsq.Producer, des task.TaskDes, currentType, currentPage int) {
	res := task.Result{
		ServerIP: "120.77.177.229",
		CityId: des.CityId,
		ErrorType: currentType,
		ErrorPage: currentPage,
		EndPage: des.PageEnd,
		Err: true,
	}

	if currentType == 74 && currentPage == des.PageEnd {
		res.Err =false
	}

	resJson, _ := json.Marshal(res)
	if err := producer.Publish("result_queue", resJson); err != nil {
		log.Fatal("publish error: " + err.Error())
	}
}

func startConsumer(topic, channel string) {
	consumer, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
	}
	consumer.AddHandler(nsq.HandlerFunc(taskHandler))
	if err := consumer.ConnectToNSQD("localhost:4150"); err != nil {
		log.Fatal(err)
	}
	<-consumer.StopChan
}

func taskHandler(message *nsq.Message) error {

	var wg sync.WaitGroup
	t := task.Task{}
	json.Unmarshal(message.Body, &t)
	fmt.Println(t)

	for i := 0; i < len(t.Describe); i++ {
		wg.Add(1)
		go func(j task.TaskDes) {
			startScrape(j)
			wg.Done()
		}(t.Describe[i])

		if (i+1)%t.Goroutines == 0 {
			wg.Wait()
		}
	}

	return nil
}

func startScrape(j task.TaskDes) {

	jobProducer, err := nsq.NewProducer("localhost:4150", nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
	}

	resultProducer, err := nsq.NewProducer("192.168.210.94:4150", nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
	}

	chrome := utils.InitClientByRemote("http://localhost:4444/wd/hub")
	defer chrome.Service.Stop()
	defer chrome.Webdriver.Quit()

	sendJobData(chrome, jobProducer, resultProducer, j)
}

func main() {
	startConsumer("task_list1", "task_channel")

	// chrome := utils.InitClientByDriver("./chromedriver", 8080, false)
	// defer chrome.Webdriver.Close()
	// defer chrome.Service.Stop()

	// sendJobData(chrome, nil, task.Job{CityId: 0, TypeId: 1, PageStart: 0, PageEnd: 1})
}
