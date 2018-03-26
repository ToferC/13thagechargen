package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/wcharczuk/go-chart"
)

type Casualty struct {
	faction      string
	stamp        time.Time
	currentForce float64
}

func (b *Combat) Report(timeTracker []timeStamp) {

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

	var casualties = make(map[string][]Casualty)

	// Preset map with starting forces value
	for k, v := range starting {
		casualties[k] = append(casualties[k],
			Casualty{
				faction:      k,
				stamp:        timeTracker[0].stamp,
				currentForce: float64(v),
			})
	}

	// Track casualties in the map
	for _, instance := range timeTracker {
		faction := instance.fighter.Faction

		count := casualties[faction][len(casualties[faction])-1].currentForce

		casualties[faction] = append(casualties[faction],
			Casualty{
				faction:      faction,
				stamp:        instance.stamp,
				currentForce: count - 1.0,
			})
	}

	serveGraph(casualties)
}

func drawChart(w http.ResponseWriter, r *http.Request, c map[string][]Casualty) {

	// Set up for multiple series of timeseries graphs
	var s []chart.Series
	var colorIndex int

	for k, v := range c {
		var xvalues []float64
		var yvalues []float64

		for _, instance := range v {
			xvalues = append(xvalues, float64(instance.stamp.Second()))
			yvalues = append(yvalues, instance.currentForce)
		}

		// Create timeseries and append to slice
		s = append(s, chart.ContinuousSeries{
			Name: k,
			Style: chart.Style{
				Show:        true,
				StrokeColor: chart.GetAlternateColor(colorIndex),
				FillColor:   chart.GetAlternateColor(colorIndex).WithAlpha(64),
			},
			XValues: xvalues,
			YValues: yvalues,
		})
		colorIndex++
	}

	graph := chart.Chart{
		YAxis: chart.YAxis{
			Style: chart.StyleShow(),
			TickStyle: chart.Style{
				TextRotationDegrees: 45.0,
			},
		},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Series: s,
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	w.Header().Set("Content-Type", "image/png")
	graph.Render(chart.PNG, w)
}

func serveGraph(c map[string][]Casualty) {
	http.HandleFunc("/favico.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte{})
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		drawChart(w, r, c)
	})
	http.ListenAndServe(":8080", nil)
}
