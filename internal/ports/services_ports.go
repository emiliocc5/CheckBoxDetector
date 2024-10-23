package ports

import (
	"image"
	"io"
)

//go:generate mockgen -destination=../mocks/services_ports_mock.go -package=mocks -source=services_ports.go
type CheckBoxDetector interface {
	Detect() (int, error)
}

type ImageGetter interface {
	GetImage(path string) (io.Reader, error)
}

type ImageDecoder interface {
	Decode(reader io.Reader) (image.Image, error)
}

type ImageGrayer interface {
	ToGray(img image.Image) *image.Gray
}

type ImageBinarizer interface {
	ApplyWithAdaptiveThreshold(img *image.Gray, blockSize, fineAdjustment int) *image.Gray
}
