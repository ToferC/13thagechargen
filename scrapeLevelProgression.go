package main

/*
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

	class := "fighter"

	path := "http://www.13thagesrd.com/classes/" + class

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

} */
