package main

import (
	"fmt"
	"math/rand"
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
	Down         bool
	Target       *Combatant
}

type Combat struct {
	fighters []*Combatant
	active   bool
}

// Check to see if the target is out of the fight
func (c *Combatant) validate() {

	if c.HP <= 0 {
		fmt.Printf("%s is struck down!\n", c.Name)
		c.Down = true
	}
}

func RollDie(max, min, numDice int) int {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	result := 0
	for i := 1; i < numDice+1; i++ {
		result += r1.Intn(max-min) + min
	}
	return result
}

func (c *Combatant) selectTarget(r *Combat) {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	var potentialTargets []*Combatant

	for _, f := range r.fighters {
		if f != c && f.Down != true && c.Faction != f.Faction {
			potentialTargets = append(potentialTargets, f)
		}
	}

	if len(potentialTargets) == 0 {
		fmt.Printf("The combat is over! The %s wins!\n", c.Faction)
		r.active = false
	} else {
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

	hugo := Combatant{Name: "Hugo", Faction: "Bloodhawks", HP: 10, AC: 15, AttackBonus: 3,
		WeaponDamage: 8, DamageBonus: 2, Initiative: 3, Down: false}

	blackthorn := Combatant{Name: "BlackThorn", Faction: "Bloodhawks", HP: 10, AC: 15, AttackBonus: 3,
		WeaponDamage: 8, DamageBonus: 2, Initiative: 3, Down: false}

	rutger := Combatant{Name: "Rutger", Faction: "Templars", HP: 13, AC: 18, AttackBonus: 4,
		WeaponDamage: 8, DamageBonus: 3, Initiative: 3, Down: false}

	r := Combat{fighters: []*Combatant{&hugo, &blackthorn, &rutger}, active: true}

	//activeCombatants := r.fighters

	for r.active {

		for _, fighter := range r.fighters {

			if fighter.Down != true {

				if r.active {

					fighter.selectTarget(&r)
				}

				if r.active {
					fighter.attack()
					fighter.Target.validate()
				}

			}
		}
	}
	fmt.Println("FIGHT OVER!")
}
