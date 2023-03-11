package main

import (
	"fmt"
	"time"
	"utils"
)

// go run scraper/*.go
func getBossJob() {
	chrome := utils.InitClientByDriver("./chromedriver", 8080, false)
	defer chrome.Webdriver.Quit()
	defer chrome.Service.Stop()

	chrome.Webdriver.Get("https://www.zhipin.com/web/geek/job?query=&city=100010000")

	time.Sleep(10 * time.Second)
	fmt.Println()
}
