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

// Check to see if the target is out of the fight
func (c *Combatant) validate(s *Combatant, tm *[]timeStamp) {

	if c.HP <= 0 {
		fmt.Printf("%s is struck down!\n", c.Name)
		c.Down = true
		ts := timeStamp{fighter: c, slayer: s, stamp: time.Now()}
		*tm = append(*tm, ts)
	}
}

// Roll and sum dice
func RollDie(max, min, numDice int) int {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	result := 0
	for i := 1; i < numDice+1; i++ {
		result += r1.Intn(max-min) + min
	}
	return result
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
		if f != c && f.Down != true && c.Faction != f.Faction {
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

func (b *Combat) report(timeTracker []timeStamp) {

	factions := make(map[string]bool)

	for _, fighter := range b.fighters {
		factions[fighter.Faction] = true
	}

	// Track starting combatants
	starting := make(map[string]int)

	for k, _ := range factions {
		for _, fighter := range b.fighters {
			if k == fighter.Faction {
				starting[k] += 1
			}
		}
	}

	// Track combatants at end of battle
	results := make(map[string]int)

	for k, _ := range factions {
		for _, fighter := range b.fighters {
			if k == fighter.Faction && fighter.Down == false {
				results[k] += 1
			}
		}
	}

	// Print out results
	for k, v := range starting {
		fmt.Println(k)
		fmt.Printf("Starting Force: %d\n", v)
		fmt.Printf("Standing: %d\n", results[k])
		fmt.Printf("Losses: %.2f%%\n\n", 100-(float64(results[k])/float64(v))*100)
	}

	for _, instance := range timeTracker {
		fmt.Printf("%s is slain by %s at %s\n", instance.fighter.Name,
			instance.slayer.Name,
			instance.stamp.Format("3:04PM"))
	}

	// develop MVP script

	/* Transform timeTracker into plottable Format

		plotData := convertTimeTracker(timeTracker, factions)

		// plotting
		p, err := plot.New()
		if err != nil {
			panic(err)
		}

		p.Title.Text = "Battle Results"
		p.X.Label.Text = "Time"
		p.Y.Label.Text = "Combatants"

		err = plotutil.AddLinePoints(p)
	}

	func convertTimeTracker(t []timeStamp, f map[string]bool) map[string][]plotXY {

		plotData := make(map[string][]plotXY)

		for k, _ := range f {
			plotData[k] := []plotXY
			for _, s := range t {
				plotData[k] = append(plotData[k], plotXY{faction: k, s.stamp)
			}
		} */
}

func main() {

	fmt.Println("*** THE BATTLE BEGINS ***\n")

	// Create initiative channel
	i := make(chan *Combatant)

	// Create WaitGroup
	wg := new(sync.WaitGroup)

	var timeTracker []timeStamp

	// Initialize combatants
	hugo := Combatant{Name: "Hugo", Faction: "Bloodhawks", HP: 10, AC: 15, AttackBonus: 3,
		WeaponDamage: 8, DamageBonus: 2, Initiative: 3, Speed: 1.0, Down: false}

	blackthorn := Combatant{Name: "BlackThorn", Faction: "Bloodhawks", HP: 10, AC: 15, AttackBonus: 3,
		WeaponDamage: 8, DamageBonus: 2, Initiative: 3, Speed: 1.4, Down: false}

	rutger := Combatant{Name: "Rutger", Faction: "Templars", HP: 13, AC: 18, AttackBonus: 4,
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
	battle.report(timeTracker)

}
