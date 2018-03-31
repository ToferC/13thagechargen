package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Combatant struct {
	Name         string
	Faction      string
	HP           int
	AC           int
	AttackBonus  int
	WeaponDamage int
	DamageBonus  int
	Initiative   int
	Speed        float64
	Down         bool
	Target       *Combatant
}

type Combat struct {
	fighters []*Combatant
	active   bool
}

type timeStamp struct {
	fighter *Combatant
	slayer  *Combatant
	stamp   time.Time
}

type plotXY struct {
	stampX time.Time
	sumY   int
}

type Faction struct {
	name    string
	allies  []string
	nemesis []string
}

// Initialize factions

var templars = Faction{
	name:   "Templars",
	allies: []string{"Heroes"},
}

var bloodhawks = Faction{
	name: "Bloodhawks",
}

var heroes = Faction{
	name:   "Heroes",
	allies: []string{"Templars"},
}

var factions = []Faction{templars, bloodhawks, heroes}

func checkAllies(c Combatant, t Combatant) bool {

	actingFaction := Faction{}

	for i, f := range factions {
		if f.name == c.Faction {
			actingFaction = factions[i]
		}
	}

	allies := actingFaction.allies

	for _, ally := range allies {
		if ally == t.Faction {
			// potential target is an ally
			return true
		}
	}
	// target is not an ally
	return false
}

// Check to see if the target is out of the fight
func (c *Combatant) validate(s *Combatant, tm *[]timeStamp) {

	if c.HP <= 0 {
		fmt.Printf("%s is struck down!\n", c.Name)
		c.Down = true
		ts := timeStamp{fighter: c, slayer: s, stamp: time.Now()}
		*tm = append(*tm, ts)
	}
}

// Insert combatants into qeue based on speed
func initiative(c *Combatant, i chan *Combatant, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(1000*c.Speed) * time.Millisecond)
	i <- c
}

// Identify and select combat targets. End if no targets available
func (c *Combatant) selectTarget(battle *Combat) {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	var potentialTargets []*Combatant

	for _, f := range battle.fighters {
		if f != c && f.Down != true && c.Faction != f.Faction && !checkAllies(*c, *f) {
			potentialTargets = append(potentialTargets, f)
		}
	}

	// If no targets, end combat and declare winner, else assign target
	if len(potentialTargets) == 0 {
		fmt.Printf("The combat is over! The %s wins!\n", c.Faction)
		battle.active = false
	} else {
		// Add decision matrix here
		c.Target = potentialTargets[r1.Intn(len(potentialTargets))]
		fmt.Printf("%s selects %s as their target.\n", c.Name, c.Target.Name)
	}
}

// Basic attack and damage roll in combat
func (c *Combatant) attack() {

	attackRoll := RollDie(20, 1, 1) + c.AttackBonus

	if attackRoll >= c.Target.AC {
		damage := RollDie(c.WeaponDamage, 1, 1) + c.DamageBonus
		c.Target.HP -= damage
		fmt.Printf("%s rolls %d and hits %s for %d damage! %s has %d HP left!\n",
			c.Name, attackRoll, c.Target.Name, damage, c.Target.Name, c.Target.HP)
	} else {
		fmt.Printf("%s attacks %s, but rolls a %d and misses!\n", c.Name,
			c.Target.Name, attackRoll)
	}
}

func main() {

	fmt.Println("*** THE BATTLE BEGINS ***\n")

	// Create initiative channel
	i := make(chan *Combatant)

	// Create WaitGroup
	wg := new(sync.WaitGroup)

	var timeTracker []timeStamp

	// Initialize combatants
	hugo := Combatant{Name: "Hugo", Faction: "Heroes", HP: 10, AC: 15, AttackBonus: 3,
		WeaponDamage: 8, DamageBonus: 2, Initiative: 3, Speed: 1.0, Down: false}

	blackthorn := Combatant{Name: "BlackThorn", Faction: "Heroes", HP: 10, AC: 15, AttackBonus: 3,
		WeaponDamage: 8, DamageBonus: 2, Initiative: 3, Speed: 1.4, Down: false}

	rutger := Combatant{Name: "Rutger", Faction: "Heroes", HP: 13, AC: 18, AttackBonus: 4,
		WeaponDamage: 12, DamageBonus: 3, Initiative: 3, Speed: 0.8, Down: false}

	battle := Combat{fighters: []*Combatant{&hugo, &blackthorn, &rutger}, active: true}

	var faction string

	for counter := 0; counter < 100; counter++ {

		switch {
		case counter < 50:
			faction = "Bloodhawks"
		default:
			faction = "Templars"
		}

		battle.fighters = append(battle.fighters, &Combatant{
			Name:         "fighter_" + faction + "_" + string(counter),
			Faction:      faction,
			HP:           5,
			AC:           11,
			AttackBonus:  1,
			WeaponDamage: 6,
			DamageBonus:  0,
			Initiative:   0,
			Speed:        1.5,
			Down:         false,
		})
	}

	// Main combat loop
	for battle.active {

		// Send fighters to initiative channel
		for _, fighter := range battle.fighters {
			wg.Add(1)
			go initiative(fighter, i, wg)
		}

		// Listen and wait for end of channel signals
		go func() {
			wg.Wait()
			close(i)
		}()

		// Listen on initiative channel and take fighter actions
		for fighter := range i {
			if fighter.Down != true {

				if battle.active {

					fighter.selectTarget(&battle)
				}

				if battle.active {
					fighter.attack()
					fighter.Target.validate(fighter, &timeTracker)
					wg.Add(1)
					go initiative(fighter, i, wg)
				}

			}
		}

	}
	fmt.Println("FIGHT OVER!\n")
	battle.Report(timeTracker)

}
