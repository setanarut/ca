package main

import (
	"ca"
	"fmt"
	"log"
	"time"

	"github.com/anthonynsimon/bild/imgio"
)

var B = ca.B
var S = ca.S

func main() {

	// otomat := ca.NewAutomaton(512, 256, 33, ca.DunesRule(32))
	// otomat := ca.NewAutomaton(512, 256, 33, CyclicRule(1, 1, 24, Moore))

	// otomat := ca.NewAutomaton(512, 256, 3, ca.LifeRule(B(3), S(2, 3)))
	otomat := ca.NewAutomaton(512, 256, 3, ca.LifeRule(B(3), S(0, 1, 2, 3, 4, 5, 6, 7, 8)))
	otomat.LifeLikeAgeEnabled = true
	// otomat.FillWithRandomStates()

	ca.Birth(otomat.Current, 100-2, 100-3)
	ca.Birth(otomat.Current, 100, 100)
	ca.Birth(otomat.Current, 102, 103)

	t := time.Now()
	for range 700 {
		otomat.Step()
	}
	fmt.Println("Simulation time taken:", time.Since(t).Seconds(), "seconds")

	err := imgio.Save("examples/demo/o.png", otomat.RenderedImage(), imgio.PNGEncoder())
	if err != nil {
		log.Fatal(err)
	}
}
