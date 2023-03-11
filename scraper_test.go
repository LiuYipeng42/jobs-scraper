package main

import (
	"fmt"

	"os"
	"strings"
	"testing"
	"utils"

	"github.com/PuerkitoBio/goquery"

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
	content, err := os.ReadFile("../jobPage.html")
	if err != nil {
		panic(err)
	}
	// (?:.|\n)

	html := string(content)

	data := utils.RegExpFindOne(html, "<div class=\"bmsg job_msg inbox\">(?:.|\n)*?</div>")
	fmt.Println(data)

}
