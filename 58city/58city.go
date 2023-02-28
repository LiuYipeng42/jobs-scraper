package main

import (
	"fmt"
	"pojo"
	"strings"
	"sync"
	"time"
	"utils"

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

	currentUrl, _ := chrome.Webdriver.CurrentURL()
	if !strings.Contains(currentUrl, "antibot") {
		return false
	}

	n := 0
	for {
		check := chrome.WaitAndFindOne("input[value=点击按钮进行验证]", 2, 1)
		if check == nil {
			return true
		}
		check.Click()
		check.Click()
		check.Click()
		time.Sleep(2000 * time.Millisecond)
		chrome.Webdriver.Get(url)
		n += 1
		if n == 5 {
			fmt.Println("anti-antibot timeout!")
			time.Sleep(5 * time.Second)
			return true
		}

	}

}

func antiBot(chrome utils.ChromeSerADri, url string) {

	antiBotClickCheck(chrome, url)

	if !antiBotLogin(chrome) {
		login(chrome)
	}

}

func getInfo(html string, url string) (job pojo.Job) {

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
	job.Postion = pos[:len(pos)-1]

	experience := utils.RegExpFindOne(html, "<span class=\"item_condition border_right_None\">.*?</span>")
	job.Experience = experience[48 : len(experience)-8]

	degree := utils.RegExpFindOne(html, "<span class=\"item_condition\">.*?</span>")
	job.Degree = degree[29 : len(degree)-7]

	temp := utils.RegExpFindAll(html, "<span class=\"pos_welfare_item\">.*?</span>")
	var tags []string
	for i := 0; i < len(temp); i++ {
		tags = append(tags, temp[i][31:len(temp[i])-7])
	}
	job.Tags = tags

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
	job.MainBusiness = strings.Split(business[strings.Index(business, ">")+1:len(business)-4], "/")

	cdes := strings.ReplaceAll(utils.RegExpFindOne(html, "<div class=\"comIntro\">.*?</p>"), " ", "")
	if len(cdes) > 0 {
		job.CDescribe = cdes[60 : len(cdes)-4]
	}

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
			// fmt.Println(linke.Text())
			url, _ := linke.GetAttribute("href")
			linke.Click()

			handle, _ := wd.WindowHandles()
			wd.SwitchWindow(handle[1])

			antiBotClickCheck(chrome, url)
			html, _ := wd.ExecuteScript("return document.documentElement.outerHTML", nil)
			job := getInfo(html.(string), url)

			wd.CloseWindow(handle[1])
			wd.SwitchWindow(handle[0])
			fmt.Print(time.Since(t), "|")

			fmt.Println(pageNum, job.Name)
		}
		pageNum++
		if pageNum <= lastPage {
			url = baseUrl + fmt.Sprintf("/pn%d", pageNum)
			wd.Get(url)
			antiBotClickCheck(chrome, url)
		} else {
			nextJobType(chrome)
			pageNum = firstPage
		}

	}

}

func nextJobType(chrome utils.ChromeSerADri) {

}

func main() {

	// chrome := utils.InitClientByRemote("http://172.17.0.2:4444")
	// wd := chrome.Webdriver
	// defer wd.Quit()

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

	for city := 0; city < 1; city++ {

		chrome := utils.InitClientByDriver("./chromedriver", 8080+city, false)
		wd := chrome.Webdriver
		defer chrome.Service.Stop()
		defer wd.Quit()

		wg.Add(1)
		go func(first, last int) {
			url := "https://" + allCitys[city] + ".58.com/zplvyoujiudian"
			wd.Get(url)
			antiBot(chrome, url)
			visitJobUrl(chrome, url, first, last)
			wg.Done()
		}(0, 1)
		time.Sleep(5 * time.Second)
	}

	wg.Wait()
}
