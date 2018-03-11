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

	// set URLs for scraping
	class := "barbarian"
	path := "http://www.13thagesrd.com/classes/" + class

	// Scrape class abilities

	// Omit non-relevant H4 nodes
	omit := []string{
		"Gold Pieces", "Navigation", " Latest products in the Open Gaming Store"}

	abilities := make(map[string]string)

	// Find h4 tags for abilities
	c.OnHTML("h4", func(e *colly.HTMLElement) {
		if stringInSlice(e.Text, omit) {
		} else {

			description := ""

			goquerySelection := e.DOM

			description += goquerySelection.NextUntilSelection(goquerySelection.Find("h4")).Text() + "\n"

			abilities[e.Text] = description
		}
	})

	// Test scraping function - class abilities table
	// Looking for span id = "Level_Progression" or table here
	c.OnHTML("th", func(e *colly.HTMLElement) {
		fmt.Println("table row", e.Text)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit(path)

	//fmt.Println(abilities)

}
