package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

// StringToLines - Convert HTML table strings into text lines
func StringToLines(s string) []string {
	var lines []string

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return lines
}

func processTable(tableObject *goquery.Selection) {
	fmt.Println("Processing table")

	// map of level(int) x column(string) x value(string)
	classMap := make(map[int]map[string]string)

	tableObject.Each(func(i int, table *goquery.Selection) {

		table.Find("tr").Each(func(rowIndex int, tr *goquery.Selection) {

			tr.Find("td").Each(func(indexOfTd int, td *goquery.Selection) {
				lines := StringToLines(td.Text())
				for elementIndex, line := range lines {

					switch elementIndex {
					case 0:
						classMap[rowIndex]["Level"] = line
					case 1:
						classMap[rowIndex]["HP"] = line
					case 2:
						classMap[rowIndex]["Feats"] = line
					case 3:
						classMap[rowIndex]["Talents"] = line
					case 4:
						classMap[rowIndex]["BLANK"] = line
					case 5:
						classMap[rowIndex]["AbilityMod"] = line
					default:
						fmt.Println("Not quite working yet")
					}

					//line = strings.TrimSpace(line)
					fmt.Printf("%s\n", line)

				}
			})
		})
	})
	fmt.Println(classMap)
}

func GetAbilities(class string) map[string]string {

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

			//fmt.Println(abilities)
		}
	})

	// Test scraping function - class abilities table
	c.OnHTML("body", func(e *colly.HTMLElement) {

		//tableRows := make(map[int][]string)
		goquerySelection := e.DOM

		// Pull the class progression table
		tableObject := goquerySelection.Find("table").Eq(3)

		fmt.Println("Found Level Progression Table")

		processTable(tableObject)

	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit(path)

	//fmt.Println(abilities)

	return abilities

}
