package main

import (
	"image"
	"image/color"
	"math/rand"
	"sync"

	"slices"

	"github.com/anthonynsimon/bild/imgio"
)

type CA struct {
	width, height int
	current, next *image.Gray
	burn          []int
	survive       []int
}

func NewCA(width, height int, burn, survive []int) *CA {
	current := image.NewGray(image.Rect(0, 0, width, height))
	next := image.NewGray(image.Rect(0, 0, width, height))

	// Rastgele başlangıç durumu
	for y := range height {
		for x := range width {
			if rand.Float32() < 0.5 {
				current.SetGray(x, y, color.Gray{Y: 255})
			} else {
				current.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

	return &CA{
		width:   width,
		height:  height,
		current: current,
		next:    next,
		burn:    burn,
		survive: survive,
	}
}

func (ca *CA) countNeighbors(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < ca.width && ny >= 0 && ny < ca.height {
				if ca.current.GrayAt(nx, ny).Y == 255 {
					count++
				}
			}
		}
	}
	return count
}

func (ca *CA) processRow(y int, wg *sync.WaitGroup) {
	defer wg.Done()

	for x := range ca.width {
		neighbors := ca.countNeighbors(x, y)
		current := ca.current.GrayAt(x, y).Y

		if current == 255 {
			survives := slices.Contains(ca.survive, neighbors)
			if survives {
				ca.next.SetGray(x, y, color.Gray{Y: 255})
			} else {
				ca.next.SetGray(x, y, color.Gray{Y: 0})
			}
		} else {
			burns := slices.Contains(ca.burn, neighbors)
			if burns {
				ca.next.SetGray(x, y, color.Gray{Y: 255})
			} else {
				ca.next.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}
}

func (ca *CA) Step() {
	var wg sync.WaitGroup

	for y := range ca.height {
		wg.Add(1)
		go ca.processRow(y, &wg)
	}

	wg.Wait()
	ca.current, ca.next = ca.next, ca.current
}

func main() {
	// Örnek kurallar:
	// burn: Ölü bir hücre bu sayıda komşuya sahipse canlanır
	// survive: Canlı bir hücre bu sayıda komşuya sahipse hayatta kalır
	burn := []int{3}       // Conway'in Yaşam Oyunu gibi
	survive := []int{2, 3} // Conway'in Yaşam Oyunu gibi

	ca := NewCA(500, 500, burn, survive)
	for range 600 {
		// print(".")
		ca.Step()
	}

	imgio.Save("a.png", ca.current, imgio.PNGEncoder())

}
