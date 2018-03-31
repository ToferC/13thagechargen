package main

import (
  "math/rand"
  "time"
)

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
