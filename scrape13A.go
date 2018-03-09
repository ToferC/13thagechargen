package main

import (
  "fmt"
  "net/http"

  "github.com/yhat/scrape"
  "golang.org/x/net/html"
  "golang.org/x/net/html/atom"
)

var class_list = map[string]int{
  "barbarian": 12,
  "fighter": 10,
  "rogue": 6,
  "wizard": 4,
}

func main() {
  // request and parse front page

  class := "barbarian"

  path := "http://www.13thagesrd.com/classes/" + class

  resp, err := http.Get(path)
  if err != nil {
    panic(err)
  }
  root, err := html.Parse(resp.Body)
  if err != nil {
    panic(err)
  }

  // set a map
  var abilityScrape = make(map[string]string)

  // grab all abilities and printCharacter
  abilities := scrape.FindAllNested(root, scrape.ByTag(atom.H4))
  for _, ability := range abilities {
    fmt.Println(scrape.Text(ability))
    abilityScrape[scrape.Text(ability)] = "default"
    for _, details := range scrape.FindAllNested(ability, scrape.ByTag(atom.P)) {
      fmt.Println(scrape.Text(ability))
      abilityScrape[scrape.Text(ability)] = scrape.Text(details)
    }
  }
  fmt.Println(abilityScrape)
}
