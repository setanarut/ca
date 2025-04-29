package cca

import (
	"fmt"
	"image"
	"math/rand/v2"
	"sync"
)

type Neighborhood int

const Neumann Neighborhood = 0
const Moore Neighborhood = 1

type Rule struct {
	Range, Threshold int
	States           uint8
	Neighborhood     Neighborhood
}

func (r *Rule) String() string {
	var neig string
	if r.Neighborhood == Moore {
		neig = "Moore"
	} else {
		neig = "Neumann"
	}
	return fmt.Sprintf("%d_%d_%d_%v", r.Range, r.Threshold, r.States, neig)
}

// Cyclic Cellular Automaton
type CCA struct {
	Grid, Buffer *image.Gray
	Neighbors    [][]int
	Rnd          *rand.Rand
	Rule         Rule
}

func NewCyclicAutomaton(w, h int, seed uint64, r Rule) *CCA {
	ca := &CCA{
		Grid:   image.NewGray(image.Rect(0, 0, w, h)),
		Buffer: image.NewGray(image.Rect(0, 0, w, h)),
		Rnd:    rand.New(rand.NewPCG(seed, 0)),
		Rule:   r,
	}
	ca.initGrid()
	ca.precomputeNeighbors()
	return ca
}

func (c *CCA) initGrid() {
	for y := range c.Grid.Rect.Dy() {
		for x := range c.Grid.Rect.Dx() {
			state := uint8(c.Rnd.IntN(int(c.Rule.States)))
			idx := c.Grid.PixOffset(x, y)
			c.Grid.Pix[idx] = state
		}
	}
}

func (c *CCA) precomputeNeighbors() {
	w, h := c.Grid.Rect.Dx(), c.Grid.Rect.Dy()
	c.Neighbors = make([][]int, len(c.Grid.Pix))
	for y := range h {
		for x := range w {
			index := c.Grid.PixOffset(x, y)
			dirs := []int{}
			for dy := -c.Rule.Range; dy <= c.Rule.Range; dy++ {
				for dx := -c.Rule.Range; dx <= c.Rule.Range; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}

					if c.Rule.Neighborhood != Moore && Abs(dx)+Abs(dy) > c.Rule.Range {
						continue
					}

					nx := (x + dx + w) % w
					ny := (y + dy + h) % h
					neighborIndex := c.Grid.PixOffset(nx, ny)
					dirs = append(dirs, neighborIndex)
				}
			}
			c.Neighbors[index] = dirs
		}
	}
}

func (c *CCA) Step() {
	w, h := c.Grid.Rect.Dx(), c.Grid.Rect.Dy()
	var wg sync.WaitGroup
	wg.Add(h)

	// Update grid to the next generation concurrently, row by row.
	for y := range h {
		// Capture current y
		go func(y int) {
			defer wg.Done()
			for x := range w {
				idx := c.Grid.PixOffset(x, y)
				currentState := c.Grid.Pix[idx]
				targetState := (currentState + 1) % c.Rule.States
				count := 0

				for _, nIdx := range c.Neighbors[idx] {
					if c.Grid.Pix[nIdx] == targetState {
						count++
						if count >= c.Rule.Threshold {
							break
						}
					}
				}

				if count >= c.Rule.Threshold {
					c.Buffer.Pix[idx] = targetState
				} else {
					c.Buffer.Pix[idx] = currentState
				}
			}
		}(y)
	}

	wg.Wait()

	// Swap grid and buffer
	c.Grid, c.Buffer = c.Buffer, c.Grid
}
