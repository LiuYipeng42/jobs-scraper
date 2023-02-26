package main

import (
	"fmt"
	"utils"
	"github.com/tebeka/selenium"
	"pojo"
)

func println(s string) {
	fmt.Println(s)
}

func checkOrLogin(chrome utils.ChromeSerADri) {
	chrome.WaitAndFindOne("input.btn_tj", 2).Click()

	login := chrome.WaitAndFindOne("button#btn_account", 2)
	if login == nil {
		return
	}
	chrome.WaitAndFindOne("input[name=username]", 2).SendKeys("18239061018")
	chrome.WaitAndFindOne("input[name=password]", 2).SendKeys("12345678bofs")
	login.Click()
}

func jobPageAntiBot (chrome utils.ChromeSerADri, url string) {
	title := chrome.WaitAndFindOne("span.pos_title", 1)
	for title == nil {
		checkOrLogin(chrome)
		title = chrome.WaitAndFindOne("span.pos_title", 2)
		if title == nil {
			chrome.Webdriver.Get(url)
		}
	}
}

func getInfo(html string) (job pojo.Job) {

	name := utils.RegExpFindOne(html, "class=\"pos_title\">.*?</span>")
	job.Name = name[19:len(name)-7]



	return
}

func visitJobUrl() {
	chrome := utils.InitDriver("./chromedriver", 8080, false)
	wd := chrome.Webdriver
	defer chrome.Service.Stop()
	defer wd.Quit()

	wd.Get("https://bj.58.com/zplvyoujiudian/")

	jobs := chrome.WaitAndFindAll("li.job_item.clearfix", 2)

	if len(jobs) == 0 {
		checkOrLogin(chrome)
	}

	jobs = chrome.WaitAndFindAll("li.job_item.clearfix", 2)

	for _, job := range jobs {
		linkE, _ := job.FindElement(selenium.ByCSSSelector, "div.job_name.clearfix>a[href]")
		url, _ := linkE.GetAttribute("href")
		linkE.Click()

		handle, _ := wd.WindowHandles()
		wd.SwitchWindow(handle[1])
		
		jobPageAntiBot(chrome, url)

		html, _ := wd.ExecuteScript("return document.documentElement.outerHTML", nil)
		getInfo(html.(string))

		wd.CloseWindow(handle[1])
		wd.SwitchWindow(handle[0])
	}

}

func main() {

	visitJobUrl()
}