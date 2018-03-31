package main

type DropdownItem struct {
	Name       string
	Value      string
	IsDisabled bool
	isChecked  bool
	Text       string
}

var ClassDropdown = map[string]interface{}{
	"Barbarian":   "barbarian",
	"Bard":        "bard",
	"Chaos Mage":  "combat mage",
	"Cleric":      "cleric",
	"Commander":   "commander",
	"Druid":       "druid",
	"Fighter":     "fighter",
	"Monk":        "monk",
	"Necromancer": "necromancer",
	"Occultist":   "occultist",
	"Paladin":     "paladin",
	"Ranger":      "ranger",
	"Rogue":       "rogue",
	"Sorcerer":    "sorceror",
	"Wizard":      "wizard",
}

var ClassStats = map[string]map[string]int{
	"Barbarian":   {"HP": 7, "AC": 12, "PD": 10, "MD": 10, "BP": 8, "Rec": 8},
	"Bard":        {"HP": 7, "AC": 12, "PD": 10, "MD": 10, "BP": 8, "Rec": 8},
	"Chaos Mage":  {"HP": 8, "AC": 12, "PD": 10, "MD": 10, "BP": 8, "Rec": 8},
	"Cleric":      {"HP": 7, "AC": 12, "PD": 10, "MD": 10, "BP": 8, "Rec": 8},
	"Commander":   {"HP": 7, "AC": 12, "PD": 10, "MD": 10, "BP": 8, "Rec": 8},
	"Druid":       {"HP": 7, "AC": 12, "PD": 10, "MD": 10, "BP": 8, "Rec": 8},
	"Fighter":     {"HP": 8, "AC": 12, "PD": 10, "MD": 10, "BP": 8, "Rec": 8},
	"Monk":        {"HP": 7, "AC": 12, "PD": 10, "MD": 10, "BP": 8, "Rec": 8},
	"Necromancer": {"HP": 6, "AC": 12, "PD": 10, "MD": 10, "BP": 8, "Rec": 8},
	"Occultist":   {"HP": 6, "AC": 12, "PD": 10, "MD": 10, "BP": 8, "Rec": 8},
	"Paladin":     {"HP": 8, "AC": 12, "PD": 10, "MD": 10, "BP": 8, "Rec": 8},
	"Ranger":      {"HP": 7, "AC": 12, "PD": 10, "MD": 10, "BP": 8, "Rec": 8},
	"Rogue":       {"HP": 6, "AC": 12, "PD": 10, "MD": 10, "BP": 8, "Rec": 8},
	"Sorcerer":    {"HP": 6, "AC": 12, "PD": 10, "MD": 10, "BP": 8, "Rec": 8},
	"Wizard":      {"HP": 6, "AC": 12, "PD": 10, "MD": 10, "BP": 8, "Rec": 8},
}

var RaceDropdown = map[string]interface{}{
	"Catfolk":     "catfolk",
	"Dark Elf":    "darkelf",
	"Dragonspawn": "dragonspawn",
	"Dwarf":       "dwarf",
	"Forgeborn":   "forgeborn",
	"Gearforged":  "gearforged",
	"Gnome":       "gnome",
	"Hagborn":     "hagborn",
	"Half-elf":    "halfelf",
	"Half-orc":    "halforc",
	"Halfling":    "halfling",
	"High Elf":    "highelf",
	"Aasimar":     "aasimar",
	"Human":       "human",
	"Lizardfolk":  "lizardfolk",
	"Merfolk":     "merfolk",
	"Monsters":    "monsters",
	"Samsarans":   "samsarans",
	"Tiefling":    "tiefling",
	"Wood Elf":    "woodelf",
}
