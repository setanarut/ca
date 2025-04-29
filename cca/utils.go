package cca

import (
	"image"
)

// Abs returns the absolute value of n.
func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func GetFrame(sim *CCA) *image.Gray {
	img := image.NewGray(sim.Grid.Rect)
	width := img.Rect.Dx()
	height := img.Rect.Dy()
	total := width * height
	for i := range total {
		x := i % width
		y := i / width
		state := sim.Grid.Pix[sim.Grid.PixOffset(x, y)]
		shade := uint8((float64(state) / float64(sim.Rule.States-1)) * 255)
		img.Pix[i] = shade
	}

	return img
}

// FindLastLongestConsecutivePattern dizinin sonuna en yakın olan ardışık tekrar eden
// en uzun desenin başlangıç indeksini, deseni ve uzunluğunu döndürür
// startIdx, pattern, length
func FindLastLongestConsecutivePattern(arr []uint8) (int, []uint8, int) {
	n := len(arr)
	if n <= 1 {
		return -1, nil, 0 // Desen bulunamadı
	}

	bestLength := 0
	bestStartIdx := -1

	// Her potansiyel başlangıç noktasını sondan başa doğru dene
	for i := n - 2; i >= 0; i-- {
		// Her potansiyel desen uzunluğunu dene (en fazla kalan dizi uzunluğunun yarısı kadar)
		for length := 1; i+2*length <= n; length++ {
			// Hemen ardından gelen desen aynı mı kontrol et
			isRepeating := true
			for k := range length {
				if arr[i+k] != arr[i+length+k] {
					isRepeating = false
					break
				}
			}

			// Eğer ardışık tekrar ediyorsa ve şu ana kadarki en uzun desenden daha uzunsa kaydet
			if isRepeating {
				// Uzunluk aynı olsa bile daha sonra gelen deseni tercih et (sondan başa doğru ilerlediğimiz için)
				if length >= bestLength {
					bestLength = length
					bestStartIdx = i
				}
				break // Bu başlangıç noktası için en uzun deseni bulduk
			}
		}
	}

	if bestStartIdx != -1 {
		return bestStartIdx, arr[bestStartIdx : bestStartIdx+bestLength], bestLength
	}
	return -1, nil, 0
}
