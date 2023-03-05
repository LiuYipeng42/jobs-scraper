package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"job"
	"log"
	"strings"
	"sync"
	"time"
	"utils"

	"github.com/nsqio/go-nsq"
	"github.com/tebeka/selenium"
)

func login(chrome utils.ChromeSerADri) {

	chrome.WaitAndFindOne("a[tongji_tag=pc_topbar_log_login]", 2, 1).Click()

	time.Sleep(1 * time.Second)

	chrome.WaitAndFindOne("input[placeholder=账号名]", 2, 1).SendKeys("18239061018")
	chrome.WaitAndFindOne("input[placeholder=密码]", 2, 1).SendKeys("12345678bofs")

	login, _ := chrome.Webdriver.FindElement(selenium.ByCSSSelector, "button[loginmode=accountlogin]")
	login.Click()
}

func antiBotLogin(chrome utils.ChromeSerADri) bool {

	check := chrome.WaitAndFindOne("input.btn_tj", 2, 1)
	if check == nil {
		return false
	}
	check.Click()
	chrome.WaitAndFindOne("input[name=username]", 2, 1).SendKeys("18239061018")
	chrome.WaitAndFindOne("input[name=password]", 2, 1).SendKeys("12345678bofs")

	login := chrome.WaitAndFindOne("button#btn_account", 2, 1)
	login.Click()

	return true
}

func antiBot(chrome utils.ChromeSerADri, url string) {

	if !antiBotLogin(chrome) {
		login(chrome)
	}

	for {
		texte := chrome.WaitAndFindOne("h1.item", 2, 1)
		if texte != nil {
			text, _ := texte.Text()
			if strings.Contains(text, "系统检测到您疑似使用网页抓取工具访问本网站") {
				time.Sleep(15 * time.Second)
			}
		} else {
			return
		}
	}

}

func parseHtml(html string) (j job.Job) {

	name := utils.RegExpFindOne(html, "class=\"pos_title\">.*?</span>")
	if len(name) == 0 {
		return
	}
	j.Name = name[19 : len(name)-7]

	// update := utils.RegExpFindOne(html, "<span class=\"pos_base_num pos_base_update\">.*?</span>")
	// job.Update = update[59 : len(update)-7]

	salary := utils.RegExpFindOne(html, "<span ((class=\"pos_salary\">.*?<span)|(class=\"pos_salary daiding\".*?)</span)")
	if strings.Contains(salary, "daiding") {
		j.Salary = salary[33 : len(salary)-6]
	} else {
		if len(salary) > 5 {
			j.Salary = salary[25 : len(salary)-5]
		}
	}

	posSlice := utils.RegExpFindAll(html, "<span class=\"pos_area_item\">.*?</span>")
	var posBuffer strings.Builder
	for i := 0; i < len(posSlice); i++ {
		posBuffer.WriteString(posSlice[i][28 : len(posSlice[i])-7])
		posBuffer.WriteString("/")
	}
	pos := posBuffer.String()
	if len(pos) > 1 {
		j.Position = pos[:len(pos)-1]
	}

	experience := utils.RegExpFindOne(html, "<span class=\"item_condition border_right_None\">.*?</span>")
	if len(experience) > 8 {
		j.Experience = experience[48 : len(experience)-8]
	}

	degree := utils.RegExpFindOne(html, "<span class=\"item_condition\">.*?</span>")
	if len(degree) > 7 {
		j.Degree = degree[29 : len(degree)-7]
	}

	tagElems := utils.RegExpFindAll(html, "<span class=\"pos_welfare_item\">.*?</span>")
	var tagsBuffer bytes.Buffer
	for i := 0; i < len(tagElems); i++ {
		tagsBuffer.WriteString(tagElems[i][31:len(tagElems[i])-7] + " ")
	}
	j.Tags = tagsBuffer.String()
	j.Tags = j.Tags[:len(j.Tags)]

	des := utils.RegExpFindOne(html, "<div class=\"des\">.*?</div>")
	des = strings.ReplaceAll(des, "<br>", "")
	des = strings.ReplaceAll(des, " ", "")
	if len(des) > 6 {
		j.Describe = des[16 : len(des)-6]
	}

	cname := utils.RegExpFindOne(html, "<div class=\"baseInfo_link\"((>.*?</div>)|( title=\".*?\">))")
	if len(cname) > 10 {
		if cname[len(cname)-10:] == "</a></div>" {
			cname = cname[:len(cname)-10]
			j.CName = cname[strings.LastIndex(cname, ">")+1:]
		} else {
			if cname[len(cname)-6:] == "</div>" {
				j.CName = cname[33 : len(cname)-6]
			} else {
				j.CName = cname[40 : len(cname)-2]
			}
		}
	}

	csize := utils.RegExpFindOne(html, "<p class=\"comp_baseInfo_scale\">.*?</p>")
	if len(csize) > 31 {
		j.CSize = csize[31 : len(csize)-4]
	}

	business := utils.RegExpFindOne(html, "<a class=\"comp_baseInfo_link\".*?</a>")
	j.MainBusiness = strings.ReplaceAll(business[strings.Index(business, ">")+1:len(business)-4], "/", " ")

	cdes := utils.RegExpFindOne(html, "<div class=\"comIntro\".*?</p>")

	if len(cdes) > 0 {
		cdes = cdes[strings.Index(cdes, "<p>")+3 : len(cdes)-4]
	}
	cdes = strings.ReplaceAll(cdes, " ", "")
	cdes = strings.ReplaceAll(cdes, "<br>", "")
	j.CDescribe = cdes

	return
}

func visitJobUrl(chrome utils.ChromeSerADri, producer *nsq.Producer, baseUrl string, firstPage, lastPage int) {

	pageNum := firstPage
	wd := chrome.Webdriver
	var url string

	wd.SetPageLoadTimeout(5000 * time.Millisecond)

	for pageNum < lastPage {
		chrome.WaitAndFindOne("a.icon_58zp", 5, 1)

		jobs := chrome.WaitAndFindAll("li.job_item.clearfix", 2)
		fmt.Print("job num: ")
		fmt.Println(len(jobs))

		for _, job := range jobs {
			linke, _ := job.FindElement(selenium.ByCSSSelector, "div.job_name.clearfix>a[href]")
			if linke == nil {
				continue
			}
			url, _ := linke.GetAttribute("href")
			linke.Click()
			// fmt.Println(linke.Text())

			handles, _ := wd.WindowHandles()
			wd.SwitchWindow(handles[1])

			html, _ := wd.ExecuteScript("return document.documentElement.outerHTML", nil)
			if html == nil {
				continue
			}
			job := parseHtml(html.(string))
			job.Url = url

			jobJson, _ := json.Marshal(job)
			if err := producer.Publish("jobs", jobJson); err != nil {
				log.Fatal("publish error: " + err.Error())
			}

			wd.CloseWindow(handles[1])
			wd.SwitchWindow(handles[0])
			fmt.Println(baseUrl, ""+fmt.Sprintf("page: %d", pageNum), "job: "+job.Name)
		}
		pageNum++
		if pageNum < lastPage {
			url = baseUrl + fmt.Sprintf("/pn%d", pageNum)
			wd.Get(url)
		}
	}

}

func allJobTypes(chrome utils.ChromeSerADri) []string {
	showjobs := chrome.WaitAndFindOne("div.crumb_item.slide_item", 2, 1)
	showjobs.Click()

	joblist, _ := chrome.Webdriver.FindElement(selenium.ByCSSSelector, "div.jobcatebox")

	typese, _ := joblist.FindElements(selenium.ByCSSSelector, "a.catename")
	types := []string{}
	for i := 0; i < len(typese); i++ {
		cln, _ := typese[i].GetAttribute("cln")
		types = append(types, cln)
	}

	showjobs.Click()

	return types
}

func start() {

	allCitys := []string{
		"bj", "sh", "gz", "sz", "cd", "hz", "nj", "tj", "wh", "cq", "hf", "wuhu", "bengbu", "fy", "hn", "anqing", "fz", "xm", "qz", "pt", "zhangzhou", "gz",
		"dg", "fs", "zs", "zh", "huizhou", "nn", "liuzhou", "gl", "yulin", "wuzhou", "bh", "gy", "zunyi", "qdn", "lz", "tianshui", "by", "qingyang", "pl",
		"haikou", "sanya", "wzs", "sansha", "qh", "zz", "luoyang", "xx", "ny", "xc", "pds", "ay", "hrb", "dq", "qqhr", "mdj", "suihua", "wh", "yc", "xf",
		"jingzhou", "shiyan", "hshi", "xiaogan", "cs", "zhuzhou", "yiyang", "changde", "hy", "xiangtan", "sjz", "bd", "ts", "lf", "hd", "qhd", "cangzhou", "su",
		"wx", "cz", "xz", "nt", "yz", "nc", "ganzhou", "jj", "yichun", "ja", "cc", "jl", "sp", "yanbian", "songyuan", "sy", "dl", "as", "jinzhou", "fushun", "yk",
		"yinchuan", "wuzhong", "hu", "bt", "chifeng", "erds", "xn", "hx", "haibei", "guoluo", "qd", "jn", "yt", "wf", "linyi", "zb", "jining", "ta", "lc", "weihai",
		"ty", "linfen", "dt", "yuncheng", "jz", "changzhi", "xa", "xianyang", "baoji", "wn", "hanzhong", "mianyang", "deyang", "nanchong", "yb", "zg", "ls", "xj",
		"changji", "bygl", "yili", "aks", "ks", "lasa", "rkz", "sn", "linzhi", "km", "qj", "dali", "honghe", "yx", "lj", "nb", "wz", "jh", "jx", "tz", "sx",
	}

	var wg sync.WaitGroup

	numOfGoroutine := 2

	startCity := 4
	stratType := 0

	firstPage := 0
	pageNumForEachJobType := 1

	chromes := make([]utils.ChromeSerADri, numOfGoroutine)

	for i := 0; i < numOfGoroutine; i++ {
		chrome := utils.InitClientByDriver("./chromedriver", 8080+i, false)
		chrome.Webdriver.ResizeWindow("", 1000, 1000)
		chromes = append(chromes, chrome)
	}

	// 遍历所有城市
	for city := startCity; city < len(allCitys); city++ {

		wg.Add(1)
		go func(chrome utils.ChromeSerADri, city, first, last int) {

			defer chrome.Service.Stop()
			defer chrome.Webdriver.Quit()

			producer, err := nsq.NewProducer("120.77.177.229:4150", nsq.NewConfig())
			if err != nil {
				log.Fatal(err)
			}

			url := "https://" + allCitys[city] + ".58.com/zplvyoujiudian"
			chrome.Webdriver.Get(url)
			antiBot(chrome, url)
			types := allJobTypes(chrome)
			fmt.Println(types[stratType], types)

			if city == startCity || city == startCity+1 {
				types = types[stratType:]
			}

			// 遍历所有工作种类
			for i := 0; i < len(types); i++ {
				fmt.Print("type: ")
				url = "https://" + allCitys[city] + ".58.com/" + types[i]
				chrome.Webdriver.Get(url)
				visitJobUrl(chrome, producer, url, first, last)
			}

			wg.Done()
		}(chromes[city%numOfGoroutine], city, firstPage, firstPage+pageNumForEachJobType)

		time.Sleep(15 * time.Second)

		if (city+1)%numOfGoroutine == 0 {
			wg.Wait()
		}
	}

}
