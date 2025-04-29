package ca

import (
	"image"
	"image/color"
	"math/rand/v2"
	"sync"

	"slices"
)

type CA struct {
	W, H          int
	Current, Next *image.Gray
	Burn          []int
	Survive       []int
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
		W:       width,
		H:       height,
		Current: current,
		Next:    next,
		Burn:    burn,
		Survive: survive,
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
			if nx >= 0 && nx < ca.W && ny >= 0 && ny < ca.H {
				if ca.Current.GrayAt(nx, ny).Y == 255 {
					count++
				}
			}
		}
	}
	return count
}

func (ca *CA) processRow(y int, wg *sync.WaitGroup) {
	defer wg.Done()

	for x := range ca.W {
		neighbors := ca.countNeighbors(x, y)
		current := ca.Current.GrayAt(x, y).Y

		if current == 255 {
			survives := slices.Contains(ca.Survive, neighbors)
			if survives {
				ca.Next.SetGray(x, y, color.Gray{Y: 255})
			} else {
				ca.Next.SetGray(x, y, color.Gray{Y: 0})
			}
		} else {
			burns := slices.Contains(ca.Burn, neighbors)
			if burns {
				ca.Next.SetGray(x, y, color.Gray{Y: 255})
			} else {
				ca.Next.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}
}

func (ca *CA) Step() {
	var wg sync.WaitGroup

	for y := range ca.H {
		wg.Add(1)
		go ca.processRow(y, &wg)
	}

	wg.Wait()
	ca.Current, ca.Next = ca.Next, ca.Current
}

func RandomFillInRect(im *image.Gray, r image.Rectangle, full bool) {
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
