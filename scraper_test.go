package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
)

func TestGoquery(t *testing.T) {
	content, err := os.ReadFile("../test.html")
	if err != nil {
		panic(err)
	}
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(content)))

	n := 1
	dom.Find("input[value=点击按钮进行验证]").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Attr("value"))
		n += 1
	})
}

// go test --run=TestRe
func TestRe(t *testing.T) {
	content, err := os.ReadFile("../jobPage.html")
	if err != nil {
		panic(err)
	}
	// (?:.|\n)

	html := string(content)

	data := utils.RegExpFindOne(html, "<span ((class=\"pos_salary\">.*?<span)|(class=\"pos_salary daiding\".*?)</span)")
	if strings.Contains(data, "daiding") {
		fmt.Println(data[33 : len(data)-6])
	} else {
		fmt.Println(data[25 : len(data)-5])
	}

}


func TestRemote(t *testing.T) {

	caps := selenium.Capabilities{"browserName": "chrome"}

	wd, _ := selenium.NewRemote(caps, "http://172.17.0.2:4444")
	wd.Get("https://baidu.com")
	html, _ := wd.ExecuteScript("return document.documentElement.outerHTML", nil)
	fmt.Println(1)
	fmt.Println(html)
	// if err != nil {

	// } else {
	// 	fmt.Println(err)
	// }

}

func Test1(t *testing.T) {

	var s interface{} = "hhhhhhhhhhhhhhh"

	fmt.Println(os.WriteFile("./test.html", []byte(s.(string)), 0666))
}