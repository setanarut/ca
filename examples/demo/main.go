package main

import (
	"ca"
	"fmt"
	"log"
	"time"

	"github.com/anthonynsimon/bild/imgio"
)

func main() {
	otomat := ca.NewAutomaton(512, 256, 33, ca.DunesRule(32))
	// otomat := NewAutomaton(512, 256, 33, CyclicRule(1, 1, 24, Moore))
	// otomat := NewAutomaton(512, 256, 33, LifeRule([]int{3}, []int{2, 3}))

	t := time.Now()
	for range 12 {
		otomat.Step()
	}
	fmt.Println("Simulation time taken:", time.Since(t).Seconds(), "seconds")

	err := imgio.Save("examples/demo/o.png", otomat.RenderedImage(), imgio.PNGEncoder())
	if err != nil {
		log.Fatal(err)
	}
}
