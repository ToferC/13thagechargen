package main

import (
  "fmt"
  "math/rand"
  "time"
  "io/ioutil"
  "net/http"
  "html/template"
  "encoding/json"
  "os"
)

// Character Object
type Character struct {
  Name string
  Stats map[string]int
  Class string
  Level int
  Race string
  HP int
}

// Roll dice
func rollDie (die int, num int) int {

  s1 := rand.NewSource(time.Now().UnixNano())
  r1 := rand.New(s1)

  result := 0
  for i := 0; i < num + 1; i++ {
    result += r1.Intn(die)
  }
  return result
}

// Determine stat modifiers
func findMod (stat int) int {

    mod := 0

    switch {
    case stat < 7:
      mod = -2
    case stat < 9:
      mod = -1
    case stat < 11:
      mod = 0
    case stat < 13:
      mod = 1
    case stat < 15:
      mod = 2
    case stat < 17:
      mod = 3
    case stat < 19:
      mod = 4
  }
  return mod
}

// Save & Load Functions
func loadCharacter(name string) (*Character, error) {
  filename := name + ".txt"
  json, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  var c []Character
  json.Unmarshal(raw, &c)
  return c
}

// JSON handlers
func (c Character) toString() string {
  return toJson(p)
}

func toJson (c interface{}) string {
  bytes, err := json.Marshal(c)
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }
  return string(bytes)
}


// View and Edit Handlers
func viewHandler(w http.ResponseWriter, r *http.Request, name string) {
  c, err := loadCharacter(name)
  if err != nil {
    http.Redirect(w, r, "/edit/"+name, http.StatusFound)
    return
  }
  renderTemplate(w, "view", c)
}


func editHandler(w http.ResponseWriter, r *http.Request, name string) {
  p, err := loadCharacter(name)
  if err != nil {
    c = &Character()
  }
  renderTemplate(w, "edit", c)
}

func renderTemplate(w http.ResponseWriter, tmpl string, c *Character) {
  err := templates.ExecuteTemplate(w, tmpl+"html", c)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}



/*
func saveHandler(w http.ResponseWriter, r *http.Request, name string) {
  json := r.Formvalue("")
}
*/


func main() {

  // Create character
  m := make(map[string]int)

  c := Character{name: "Baxor", stats: m, class: "",
    level: 1, race: "", hp: 0}

  c.stats["STR"] = rollDie(6,3)
  c.stats["DEX"] = rollDie(6,3)
  c.stats["CON"] = rollDie(6,3)
  c.stats["INT"] = rollDie(6,3)
  c.stats["WIS"] = rollDie(6,3)
  c.stats["CHA"] = rollDie(6,3)

  c.class = "Fighter"
  c.race = "Elf"

  conMod := 0

  // Figure out stat mod
  conMod = findMod(c.stats["CON"])

  c.hp = rollDie(10,1) + conMod

  // Print test for character
  fmt.Println("")
  fmt.Println(c.name)
  for key, value := range c.stats {
    fmt.Println(key, value)
  }
  fmt.Printf("ConMod: %d\n\n", conMod)
  fmt.Printf("%v\n\n", c)

  save := toJson(c)

  fmt.Println(save)

  // Web App
  http.HandleFunc("/view/", makeHandler(viewHandler))
  http.HandleFunc("/edit/", makeHandler(editHandler))

  log.Fatal(http.ListenAndServe(":8080", nil))


}
