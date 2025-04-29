package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand/v2"
	"sync"
	"time"

	"slices"

	"github.com/anthonynsimon/bild/imgio"
)

type CA struct {
	width, height int
	current, next *image.Gray
	burn          []int
	survive       []int
}

// NewCA kurallar
//
// burn: Ölü bir hücre bu sayıda komşulara sahipse canlanır
//
// survive: Canlı bir hücre bu sayıda komşulara sahipse hayatta kalır
func NewCA(width, height int, burn, survive []int) *CA {
	current := image.NewGray(image.Rect(0, 0, width, height))
	next := image.NewGray(image.Rect(0, 0, width, height))
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

	// B3/S012345678	Life without Death
	// B3/23	Conway

	burn := []int{3}
	survive := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}

	ca := NewCA(800, 512, burn, survive)
	x0 := (ca.width - 10) / 2
	y0 := (ca.height - 10) / 2
	randomFillInRect(ca.current, image.Rect(x0, y0, x0+10, y0+10), false)

	imgio.Save("output/ilkKare.png", ca.current, imgio.PNGEncoder())

	now := time.Now()
	// 600 adım at
	for range 1200 {
		ca.Step()
	}
	fmt.Println("Elapsed time:", time.Since(now).Seconds())

	imgio.Save("output/sonKare.png", ca.current, imgio.PNGEncoder())

}

func randomFillInRect(im *image.Gray, r image.Rectangle, full bool) {
	if full {
		r = im.Bounds()
	}
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			if rand.IntN(2) == 0 {
				im.SetGray(x, y, color.Gray{Y: 255})
			} else {
				im.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

}
