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
	content, err := os.ReadFile("../test.html")
	if err != nil {
		panic(err)
	}
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(content)))

	dom.Find("div.content-province").Each(func(i int, province *goquery.Selection) {

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

	data := utils.RegExpFindAll(html, "<a class=\"catename\" href=\".*?\" cln=\".*?\">")
	fmt.Println(data)

}