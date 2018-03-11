package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

var class_list = map[string]int{
	"barbarian": 12,
	"fighter":   10,
	"rogue":     6,
	"wizard":    4,
}

type Class struct {
	Name         string
	Melee        string
	Range        string
	Armour       []string
	HP           string
	AbilityPicks map[int]map[string]int
	Abilities    map[string]string
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {

	class := "barbarian"
	path := "http://www.13thagesrd.com/classes/" + class
	omit := []string{
		"Gold Pieces", "Navigation", " Latest products in the Open Gaming Store"}

	abilities := make(map[string]string)

	// Initialize Colly Collector
	c := colly.NewCollector(
		colly.AllowedDomains("www.13thagesrd.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	/* Scraping and traversing
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	}) */

	c.OnHTML("h4", func(e *colly.HTMLElement) {
		if stringInSlice(e.Text, omit) {
		} else {
			fmt.Println("Found Ability:", e.Text)
			description := e.ChildText("p")
			abilities[e.Text] = description
		}
	})

	/* Test scraping function - table
	c.OnHTML("tr td", func(e *colly.HTMLElement) {
		fmt.Println("table row", e.Text)
	}) */

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit(path)

	fmt.Println(abilities)

}
