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
	content, err := os.ReadFile("test.html")
	if err != nil {
		panic(err)
	}
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(content)))

	n := 1
	dom.Find("span.d.at>span").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(n, selection.Text())
		n += 1
	})
}

//  go test --run=TestRe
func TestRe(t *testing.T) {
	content, err := os.ReadFile("test.html")
	if err != nil {
		panic(err)
	}
	// (?:.|\n)

	html := string(content)

	data := utils.RegExpFindOne(html, "class=\"pos_title\">.*?</span>")
	fmt.Println(data[19:len(data)-7])

	// for i, data := range jobInfo {
	// 	fmt.Println(i, data)
	// }

}
