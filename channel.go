package main

import (
  "fmt"
  "time"
  "sync"
)

func initiative (k string, v float64, i chan string, wg *sync.WaitGroup) {
  defer wg.Done()
  time.Sleep(time.Duration(1000*v)*time.Millisecond)
  i <- k
}

func main() {
  wg := new(sync.WaitGroup)
  i := make(chan string)

  combatants := map[string]float64{"Hugo": 1.0, "Rutger": 3.5, "BlackRock": 0.7}

  for k, v := range combatants {
    wg.Add(1)
    go initiative(k, v, i, wg)
  }

  go func() {
    wg.Wait()
    close(i)
  }()

  for k := range i {
    fmt.Printf("%s acts!\n\n", k)
    v := combatants[k]
    wg.Add(1)
    go initiative(k, v, i, wg)
  }

}
