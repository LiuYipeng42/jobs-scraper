package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"pojo"
	"strings"
	"sync"
	"time"
	"utils"

	"github.com/nsqio/go-nsq"
	"github.com/tebeka/selenium"
)

func login(chrome utils.ChromeSerADri) {

	chrome.WaitAndFindOne("a[tongji_tag=pc_topbar_log_login]", 2, 1).Click()

	time.Sleep(2 * time.Second)

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

func antiBotClickCheck(chrome utils.ChromeSerADri, url string) bool {

	check := chrome.WaitAndFindOne("input[value=点击按钮进行验证]", 1, 0.5)

	if check == nil {
		return false
	}

	n := 0
	for {
		check := chrome.WaitAndFindOne("input[value=点击按钮进行验证]", 2, 1)
		if check == nil {
			return true
		}
		time.Sleep(1000 * time.Millisecond)
		check.Click()
		time.Sleep(1500 * time.Millisecond)
		chrome.Webdriver.Get(url)
		n += 1
		if n == 5 {
			fmt.Println("anti-antibot timeout!")
			time.Sleep(5 * time.Second)
		}
	}

}

func antiBot(chrome utils.ChromeSerADri, url string) {

	// https://callback.58.com/antibot/deny.do?namespace=zhaopin_list_pc&serialID=9fc606871dd48050911342ddd1cf928d_62e113b8e3714a28be99e2e47f23cc1d

	antiBotClickCheck(chrome, url)

	if !antiBotLogin(chrome) {
		login(chrome)
	}

	texte := chrome.WaitAndFindOne("h1.item", 2, 1)
	if texte != nil {
		text, _ := texte.Text()
		if strings.Contains(text, "系统检测到您疑似使用网页抓取工具访问本网站") {
			time.Sleep(10 * time.Second)
			chrome.Webdriver.Get(url)
		}	
	}

}

func getInfo(html string) (job pojo.Job) {

	name := utils.RegExpFindOne(html, "class=\"pos_title\">.*?</span>")
	if len(name) == 0 {
		return pojo.Job{}
	}
	job.Name = name[19 : len(name)-7]

	// update := utils.RegExpFindOne(html, "<span class=\"pos_base_num pos_base_update\">.*?</span>")
	// job.Update = update[59 : len(update)-7]

	salary := utils.RegExpFindOne(html, "<span ((class=\"pos_salary\">.*?<span)|(class=\"pos_salary daiding\".*?)</span)")
	if strings.Contains(salary, "daiding") {
		job.Salary = salary[33 : len(salary)-6]
	} else {
		job.Salary = salary[25 : len(salary)-5]
	}

	posSlice := utils.RegExpFindAll(html, "<span class=\"pos_area_item\">.*?</span>")
	var posBuffer strings.Builder
	for i := 0; i < len(posSlice); i++ {
		posBuffer.WriteString(posSlice[i][28 : len(posSlice[i])-7])
		posBuffer.WriteString("/")
	}
	pos := posBuffer.String()
	job.Position = pos[:len(pos)-1]

	experience := utils.RegExpFindOne(html, "<span class=\"item_condition border_right_None\">.*?</span>")
	job.Experience = experience[48 : len(experience)-8]

	degree := utils.RegExpFindOne(html, "<span class=\"item_condition\">.*?</span>")
	job.Degree = degree[29 : len(degree)-7]

	tagElems := utils.RegExpFindAll(html, "<span class=\"pos_welfare_item\">.*?</span>")
	var tagsBuffer bytes.Buffer
	for i := 0; i < len(tagElems); i++ {
		tagsBuffer.WriteString(tagElems[i][31:len(tagElems[i])-7] + " ")
	}
	job.Tags = tagsBuffer.String()
	job.Tags = job.Tags[:len(job.Tags)]

	des := utils.RegExpFindOne(html, "<div class=\"des\">.*?</div>")
	des = strings.ReplaceAll(des, "<br>", "")
	des = strings.ReplaceAll(des, " ", "")
	job.Describe = des[16 : len(des)-6]

	cname := utils.RegExpFindOne(html, "<div class=\"baseInfo_link\"((>.*?</div>)|( title=\".*?\">))")

	if cname[len(cname)-6:] == "</div>" {
		cname = cname[:len(cname)-10]
		job.CName = cname[strings.LastIndex(cname, ">")+1:]
	} else {
		job.CName = cname[34 : len(cname)-2]
	}

	csize := utils.RegExpFindOne(html, "<p class=\"comp_baseInfo_scale\">.*?</p>")
	job.CSize = csize[31 : len(csize)-4]

	business := utils.RegExpFindOne(html, "<a class=\"comp_baseInfo_link\".*?</a>")
	job.MainBusiness = strings.ReplaceAll(business[strings.Index(business, ">")+1:len(business)-4], "/", " ")

	cdes := utils.RegExpFindOne(html, "<div class=\"comIntro\".*?</p>")

	if len(cdes) > 0 {
		cdes = cdes[strings.Index(cdes, "<p>") + 3:len(cdes) - 4]
	}
	cdes = strings.ReplaceAll(cdes, " ", "")
	cdes = strings.ReplaceAll(cdes, "<br>", "")
	job.CDescribe = cdes

	return
}

func visitJobUrl(chrome utils.ChromeSerADri, baseUrl string, firstPage, lastPage int) {

	pageNum := firstPage
	wd := chrome.Webdriver
	var url string

	for pageNum < lastPage {
		chrome.WaitAndFindOne("a.icon_58zp", 5, 1)

		jobs := chrome.WaitAndFindAll("li.job_item.clearfix", 2)
		fmt.Println(len(jobs))

		for _, job := range jobs {
			t := time.Now()
			linke, _ := job.FindElement(selenium.ByCSSSelector, "div.job_name.clearfix>a[href]")
			url, _ := linke.GetAttribute("href")
			linke.Click()
			// fmt.Println(linke.Text())

			handles, _ := wd.WindowHandles()
			wd.SwitchWindow(handles[1])

			antiBotClickCheck(chrome, url)

			html, _ := wd.ExecuteScript("return document.documentElement.outerHTML", nil)
			job := getInfo(html.(string))
			job.Url = url

			jobJson, _ := json.Marshal(job)
			if err := producer.Publish("jobs", jobJson); err != nil {
				log.Fatal("publish error: " + err.Error())
			}

			wd.CloseWindow(handles[1])
			wd.SwitchWindow(handles[0])
			fmt.Println(time.Since(t), "|", pageNum, job.Name)
		}
		pageNum++
		if pageNum <= lastPage {
			url = baseUrl + fmt.Sprintf("/pn%d", pageNum)
			wd.Get(url)
			antiBotClickCheck(chrome, url)
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

	firstPage := 0
	pageNumForEachJobType := 1

	// 遍历所有城市
	for city := 0; city < len(allCitys); city++ {

		chrome := utils.InitClientByDriver("./chromedriver", 8080+city, false)
		wd := chrome.Webdriver
		defer chrome.Service.Stop()
		defer wd.Quit()

		wg.Add(1)
		go func(first, last int) {
			url := "https://" + allCitys[city] + ".58.com/zplvyoujiudian"
			wd.Get(url)
			antiBot(chrome, url)
			types := allJobTypes(chrome)

			// 遍历所有工作种类
			for i := 1; i < len(types); i++ {
				visitJobUrl(chrome, url, first, last)
				// time.Sleep(2 * time.Second)
				url = "https://" + allCitys[city] + ".58.com/" + types[i]
				wd.Get(url)
				fmt.Println(url)
				antiBotClickCheck(chrome, url)
			}

			wg.Done()
		}(firstPage, firstPage+pageNumForEachJobType)

		time.Sleep(10 * time.Second)

		if city%2 == 0 {
			wg.Wait()
		}
	}

}
