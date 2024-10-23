package adapters

import (
	"image"
	"image/color"
)

const (
	eightConstant = 8
)

type InternalImageGrayer struct{}

func NewInternalImageGrayer() *InternalImageGrayer {
	return &InternalImageGrayer{}
}

// TODO Add nil defensive clause
func (i *InternalImageGrayer) ToGray(img image.Image) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			r, g, b, _ := originalColor.RGBA()
			grayColor := uint8((r*299 + g*587 + b*114) / 1000 >> eightConstant)
			gray.Set(x, y, color.Gray{Y: grayColor})
		}
	}

	return gray
}
