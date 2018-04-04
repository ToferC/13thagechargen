package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Character Object
type Character struct {
	Name      string            "json:'name'"
	Stats     map[string]int    "json:'stats'"
	Class     string            "json:'class'"
	Level     int               "json:'level'"
	Race      string            "json:'race'"
	HP        int               "json:'hp'"
	Abilities map[string]string "json:'abilities'"
}

const baseTemplate = "templates/layout.html"

// Determine stat modifiers
func FindMod(stat int) int {

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
	fmt.Println("***" + c.Name + "***")
	fmt.Printf("Level %d %s %s \n", c.Level, c.Race, c.Class)
	for key, value := range c.Stats {
		fmt.Println(key, value, FindMod(value))
	}

	conMod := FindMod(c.Stats["CON"])

	fmt.Printf("ConMod: %d\n", conMod)
	fmt.Printf("HP: %d \n\n", c.HP)
	fmt.Printf("%v\n\n", c)
}

func toJson(c interface{}) string {
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

	render(w, "templates/index.html", charNames)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Path[len("/view/"):]

	path := "./characters/" + name + ".json"
	c := openCharacter(path)

	render(w, "templates/character.html", c)
}

func newCharHandler(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Path[len("/new/"):]

	// Create character
	m := make(map[string]int)

	c := Character{Name: name, Stats: m, Class: "",
		Level: 1, Race: "", HP: 0}

	if r.Method == "GET" {

		c.Stats["STR"] = RollDie(6, 1, 3)
		c.Stats["DEX"] = RollDie(6, 1, 3)
		c.Stats["CON"] = RollDie(6, 1, 3)
		c.Stats["INT"] = RollDie(6, 1, 3)
		c.Stats["WIS"] = RollDie(6, 1, 3)
		c.Stats["CHA"] = RollDie(6, 1, 3)

		c.Class = "Fighter"
		c.Race = "Elf"

		var conMod int

		// Figure out stat mod
		conMod = FindMod(c.Stats["CON"])

		c.HP = RollDie(10, 1, 1) + conMod

		render(w, "templates/new_char.html", c)

	} else {

		c := Character{}

		m := make(map[string]int)

		c.Stats = m

		c.Name = r.FormValue("name")
		c.Class = r.FormValue("class")
		c.Race = r.FormValue("race")
		c.Level, _ = strconv.Atoi(r.FormValue("level"))
		c.HP, _ = strconv.Atoi(r.FormValue("hp"))
		c.Stats["STR"], _ = strconv.Atoi(r.FormValue("STR"))
		c.Stats["DEX"], _ = strconv.Atoi(r.FormValue("DEX"))
		c.Stats["CON"], _ = strconv.Atoi(r.FormValue("CON"))
		c.Stats["INT"], _ = strconv.Atoi(r.FormValue("INT"))
		c.Stats["WIS"], _ = strconv.Atoi(r.FormValue("WIS"))
		c.Stats["CHA"], _ = strconv.Atoi(r.FormValue("CHA"))

		c.Abilities = GetAbilities(c.Class)

		fmt.Println(c)
		c.save()
		http.Redirect(w, r, "/view/"+c.Name, http.StatusSeeOther)
	}
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

	// Start battle?
	createReader = bufio.NewReader(os.Stdin)
	fmt.Print("Start a battle (Y/N): ")
	startBattle, _ := createReader.ReadString('\n')

	battleOrNot := strings.Trim(startBattle, " \n")

	if battleOrNot == "Y" {
		Battle()
	} else {
		fmt.Println("No fighting today.")
	}

	fmt.Println("Starting Web server at port 8080")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/new/", newCharHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))

	// print env
	env := os.Getenv("APP_ENV")
	if env == "production" {
		log.Println("Running api server in production mode")
	} else {
		log.Println("Running api server in dev mode")
	}
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

	c.Stats["STR"] = RollDie(6, 1, 3)
	c.Stats["DEX"] = RollDie(6, 1, 3)
	c.Stats["CON"] = RollDie(6, 1, 3)
	c.Stats["INT"] = RollDie(6, 1, 3)
	c.Stats["WIS"] = RollDie(6, 1, 3)
	c.Stats["CHA"] = RollDie(6, 1, 3)

	c.Class = "Fighter"
	c.Race = "Elf"

	conMod := 0

	// Figure out stat mod
	conMod = FindMod(c.Stats["CON"])

	c.HP = RollDie(10, 1, 1) + conMod

	c.save()
}

func (c Character) save() {

	printCharacter(c)

	path := "./characters/" + c.Name + ".json"

	writeFile(path, c)
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
	fmt.Println("==> done creating file", path+"\n")

	characterJson, _ := json.Marshal(c)
	err = ioutil.WriteFile(path, characterJson, 0644)
}

// Render HTML templates
func render(w http.ResponseWriter, filename string, data interface{}) {

	tmpl := make(map[string]*template.Template)

	tmpl[filename] = template.Must(template.ParseFiles(filename, baseTemplate))

	if err := tmpl[filename].ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func openCharacter(path string) Character {

	jsonFile, err := os.Open(path)
	checkError(err)
	fmt.Println("Successfully opened file: " + path + " \n")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var c Character

	json.Unmarshal(byteValue, &c)

	return c
}
