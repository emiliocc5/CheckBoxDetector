package adapters

import (
	"image"
	"image/color"
)

const (
	whitePixelValue = 255
	blackPixelValue = 0
	twoConstant     = 2
	zeroConstant    = 0
)

type InternalImageBinarizer struct{}

func NewInternalImageBinarizer() *InternalImageBinarizer {
	return &InternalImageBinarizer{}
}

// TODO Add nil defensive clause.
func (i *InternalImageBinarizer) ApplyWithAdaptiveThreshold(
	img *image.Gray, blockSize, fineAdjustment int,
) *image.Gray {
	bounds := img.Bounds()
	bin := image.NewGray(bounds)

	for yAxis := blockSize / twoConstant; yAxis < bounds.Max.Y-blockSize/twoConstant; yAxis++ {
		for xAxis := blockSize / twoConstant; xAxis < bounds.Max.X-blockSize/twoConstant; xAxis++ {
			// Calculate the surrounding media
			mean := calcMean(img, xAxis, yAxis, blockSize)
			if img.GrayAt(xAxis, yAxis).Y < uint8(mean-fineAdjustment) {
				bin.SetGray(xAxis, yAxis, color.Gray{Y: blackPixelValue})
			} else {
				bin.SetGray(xAxis, yAxis, color.Gray{Y: whitePixelValue})
			}
		}
	}

	return bin
}

// Calc mean in "neighborhood" pixels.
func calcMean(img *image.Gray, xCoord, yCoord, blockSize int) int {
	sum := zeroConstant
	count := zeroConstant
	halfSize := blockSize / twoConstant

	for dy := -halfSize; dy <= halfSize; dy++ {
		for dx := -halfSize; dx <= halfSize; dx++ {
			sum += int(img.GrayAt(xCoord+dx, yCoord+dy).Y)
			count++
		}
	}

	return sum / count
}
