package main

import (
  "fmt"
  "math/rand"
  "time"
  "encoding/json"
  "os"
  "io/ioutil"
  "bufio"
  "strings"
  "log"
  "net/http"
  "html/template"
  "path/filepath"
)

// Character Object
type Character struct {
  Name string           "json:'name'"
  Stats map[string]int  "json:'stats'"
  Class string          "json:'class'"
  Level int             "json:'level'"
  Race string           "json:'race'"
  HP int                "json:'hp'"
}

// Roll dice
func rollDie (max, min, numDice int) int {

  s1 := rand.NewSource(time.Now().UnixNano())
  r1 := rand.New(s1)

  result := 0
  for i := 0; i < numDice + 1; i++ {
    result += r1.Intn(max - min) + min
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

// JSON handlers
func (c Character) toString() string {
  return toJson(c)
}

func printCharacter(c Character) {
  // Print test for character
  fmt.Println("")
  fmt.Println("***"+c.Name+"***")
  fmt.Printf("Level %d %s %s \n", c.Level, c.Race, c.Class)
  for key, value := range c.Stats {
    fmt.Println(key, value, findMod(value))
  }

  conMod := findMod(c.Stats["CON"])

  fmt.Printf("ConMod: %d\n", conMod)
  fmt.Printf("HP: %d \n\n", c.HP)
  fmt.Printf("%v\n\n", c)
}

func toJson (c interface{}) string {
  bytes, err := json.Marshal(c)
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }
  return string(bytes)
}

func checkError(err error) {
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(0)
  }
}

// View and Edit Handlers
func indexHandler(w http.ResponseWriter, r *http.Request) {

  files, err := filepath.Glob("./characters/*")
  if err != nil {
    log.Fatal(err)
  }

  var charNames []string

  for _, path := range files {
    _, json := filepath.Split(path)
    s := strings.Split(json, ".")
    charNames = append(charNames, s[0])
  }

  t, _ := template.ParseFiles("index.html")
  t.Execute(w, charNames)

}

func viewHandler(w http.ResponseWriter, r *http.Request) {

  name := r.URL.Path[len("/view/"):]

  path := "./characters/"+name+".json"
  c := openCharacter(path)

  json, err := json.Marshal(c)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(json)
}

func main() {
  // New character?
  createReader := bufio.NewReader(os.Stdin)
  fmt.Print("Create a new character (Y/N): ")
  createCharacter, _ := createReader.ReadString('\n')

  createOrNot := strings.Trim(createCharacter, " \n")

  if createOrNot == "Y" {
      createChar()
  } else {
    fmt.Println("No character created.")
  }

  fmt.Println("Starting Web server at port 8080")
  http.HandleFunc("/", indexHandler)
  http.HandleFunc("/view/", viewHandler)

  log.Fatal(http.ListenAndServe(":8080", nil))
}


func createChar() {
  // Get user input
  reader := bufio.NewReader(os.Stdin)
  fmt.Print("Enter character name: ")
  text, _ := reader.ReadString('\n')

  characterName := strings.Trim(text, " \n")

  // Create character
  m := make(map[string]int)

  c := Character{Name: characterName, Stats: m, Class: "",
    Level: 1, Race: "", HP: 0}

  c.Stats["STR"] = rollDie(6, 1, 3)
  c.Stats["DEX"] = rollDie(6, 1, 3)
  c.Stats["CON"] = rollDie(6, 1, 3)
  c.Stats["INT"] = rollDie(6, 1, 3)
  c.Stats["WIS"] = rollDie(6, 1, 3)
  c.Stats["CHA"] = rollDie(6, 1, 3)

  c.Class = "Fighter"
  c.Race = "Elf"

  conMod := 0

  // Figure out stat mod
  conMod = findMod(c.Stats["CON"])

  c.HP = rollDie(10,1, 1) + conMod

  printCharacter(c)

  path := "./characters/"+c.Name+".json"

  writeFile(path, c)

  d := openCharacter(path)

  printCharacter(d)
}

func writeFile(path string, c Character) {

  // Check if file exists
  var _, err = os.Stat(path)

  // Create new file if needed
  if os.IsNotExist(err) {
    var file, err = os.Create(path)
    checkError(err)
    defer file.Close()
  }
  fmt.Println("==> done creating file", path +"\n")

  characterJson, _ := json.Marshal(c)
  err = ioutil.WriteFile(path, characterJson, 0644)
}

func openCharacter(path string) Character {

  jsonFile, err := os.Open(path)
  checkError(err)
  fmt.Println("Successfully opened file: "+path+" \n")

  defer jsonFile.Close()

  byteValue, _ := ioutil.ReadAll(jsonFile)

  var c Character

  json.Unmarshal(byteValue, &c)

  return c
}
