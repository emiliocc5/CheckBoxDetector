package adapters

import (
	"image"
	"image/color"
)

type InternalImagePixelHandler struct{}

func NewInternalImagePixelHandler() *InternalImagePixelHandler {
	return &InternalImagePixelHandler{}
}

// TODO Add nil defensive clause
func (i *InternalImagePixelHandler) Dilate(grayImg image.Gray, iterations int) *image.Gray {
	bounds := grayImg.Bounds()
	dilated := image.NewGray(bounds)

	for range iterations {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				neighbors := getNeighborhood(&grayImg, x, y)
				maxValue := getMaxGray(neighbors)
				dilated.SetGray(x, y, color.Gray{Y: maxValue})
			}
		}
	}

	return dilated
}

// Obtener los píxeles vecinos de una posición (x, y)
func getNeighborhood(img *image.Gray, x, y int) []uint8 {
	var neighbors []uint8
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if x+dx >= 0 && x+dx < img.Bounds().Max.X && y+dy >= 0 && y+dy < img.Bounds().Max.Y {
				neighbors = append(neighbors, img.GrayAt(x+dx, y+dy).Y)
			}
		}
	}

	return neighbors
}

// Obtener el valor máximo de una lista de píxeles en escala de grises
func getMaxGray(neighbors []uint8) uint8 {
	maxPxValue := uint8(0)
	for _, val := range neighbors {
		if val > maxPxValue {
			maxPxValue = val
		}
	}

	return maxPxValue
}
