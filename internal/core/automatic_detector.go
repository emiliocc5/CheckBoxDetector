package core

import (
	"errors"
	"fmt"
	"image"
	"log/slog"

	"github.com/emiliocc5/CheckBoxDetector/config"
	"github.com/emiliocc5/CheckBoxDetector/internal/ports"
)

const (
	iterations          = 1
	blockSize           = 20
	fineAdjustment      = 1
	minSize             = 400
	maxSize             = 600
	minAspect           = 0.8
	maxAspect           = 1.2
	minPercentageFilled = 0.4
)

type AutomaticDetector struct {
	config         config.ServiceConfiguration
	imageGetter    ports.ImageGetter
	imageDecoder   ports.ImageDecoder
	imageGrayer    ports.ImageGrayer
	imageBinarizer ports.ImageBinarizer
}

func NewAutomaticDetector(
	config config.ServiceConfiguration,
	imageGetter ports.ImageGetter,
	imageDecoder ports.ImageDecoder,
	imageGrayer ports.ImageGrayer,
	imageBinarizer ports.ImageBinarizer,
) *AutomaticDetector {
	return &AutomaticDetector{
		config:         config,
		imageGetter:    imageGetter,
		imageDecoder:   imageDecoder,
		imageGrayer:    imageGrayer,
		imageBinarizer: imageBinarizer,
	}
}

func (d *AutomaticDetector) Detect() (int, error) {
	log := slog.Default()
	var filledCheckboxes []image.Rectangle

	// Get Image
	file, errGettingImage := d.imageGetter.GetImage(d.config.File.Path)
	if errGettingImage != nil {
		log.Error("Error getting image", slog.String("error", errGettingImage.Error()))

		return 0, errors.New("error getting image")
	}

	// Decode image
	decodedImg, errDecoding := d.imageDecoder.Decode(file)
	if errDecoding != nil {
		log.Error("Error decoding image", slog.String("error", errDecoding.Error()))

		return 0, errors.New("error decoding image")
	}

	// Convert image to gray scale
	grayImage := d.imageGrayer.ToGray(decodedImg)

	// Binarize the image using an adaptive threshold
	binarizedImage := d.imageBinarizer.ApplyWithAdaptiveThreshold(grayImage, blockSize, fineAdjustment)

	// Detect contours
	contours := detectContours(binarizedImage)

	// Filter rectangles that match with filled checkbox style
	filteredRectangles := filterRectangles(contours, minSize, maxSize, minAspect, maxAspect)

	// Print filled checkboxes
	for _, rect := range filteredRectangles {
		if isCheckboxMarked(binarizedImage, rect) {
			log.Info(fmt.Sprintf(
				"Filled checkbox with initial point at coordinates: %v and final point at coordinates: %v\n",
				rect.Min, rect.Max))
			filledCheckboxes = append(filledCheckboxes, rect)
		}
	}

	return len(filledCheckboxes), nil
}

func detectContours(bin *image.Gray) []image.Rectangle {
	bounds := bin.Bounds()
	visited := make([][]bool, bounds.Max.Y)
	for i := range visited {
		visited[i] = make([]bool, bounds.Max.X)
	}

	var contours []image.Rectangle

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if bin.GrayAt(x, y).Y == 0 && !visited[y][x] {
				rect := traceContour(bin, x, y, visited)
				if rect != image.Rect(0, 0, 0, 0) {
					contours = append(contours, rect)
				}
			}
		}
	}

	return contours
}

// TODO This method could be improved
// Trace contour & return it.
func traceContour(img *image.Gray, startX, startY int, visited [][]bool) image.Rectangle { //nolint: cyclop
	bounds := img.Bounds()
	minX, minY, maxX, maxY := startX, startY, startX, startY

	stack := []image.Point{{X: startX, Y: startY}}

	for len(stack) > 0 {
		point := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		xPoint, yPoint := point.X, point.Y
		if xPoint < bounds.Min.X || xPoint >= bounds.Max.X || yPoint < bounds.Min.Y || yPoint >= bounds.Max.Y {
			continue
		}
		if visited[yPoint][xPoint] || img.GrayAt(xPoint, yPoint).Y != 0 {
			continue
		}

		visited[yPoint][xPoint] = true

		// Actualize contour limits
		if xPoint < minX {
			minX = xPoint
		}
		if xPoint > maxX {
			maxX = xPoint
		}
		if yPoint < minY {
			minY = yPoint
		}
		if yPoint > maxY {
			maxY = yPoint
		}

		stack = append(stack, image.Point{X: xPoint + 1, Y: yPoint})
		stack = append(stack, image.Point{X: xPoint - 1, Y: yPoint})
		stack = append(stack, image.Point{X: xPoint, Y: yPoint + 1})
		stack = append(stack, image.Point{X: xPoint, Y: yPoint - 1})
	}

	rect := image.Rect(minX, minY, maxX, maxY)

	width, height := rect.Dx(), rect.Dy()
	if width > 10 && height > 10 {
		return rect
	}

	return image.Rectangle{}
}

func filterRectangles(rectangles []image.Rectangle, minSize,
	maxSize int, minAspect, maxAspect float64,
) []image.Rectangle {
	filtered := make([]image.Rectangle, 0)

	for _, rect := range rectangles {
		width := rect.Dx()
		height := rect.Dy()
		area := width * height

		if area < minSize || area > maxSize {
			continue
		}

		aspectRatio := float64(width) / float64(height)
		if aspectRatio < minAspect || aspectRatio > maxAspect {
			continue
		}

		filtered = append(filtered, rect)
	}

	return filtered
}

func isCheckboxMarked(img *image.Gray, rect image.Rectangle) bool {
	blackPixelCount := 0
	totalPixelCount := rect.Dx() * rect.Dy()

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			if img.GrayAt(x, y).Y == 0 {
				blackPixelCount++
			}
		}
	}

	percentageFilled := float64(blackPixelCount) / float64(totalPixelCount)

	return percentageFilled > minPercentageFilled
}
