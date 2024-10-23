package adapters

import (
	"image"
	"image/jpeg"
	"io"
)

type InternalImageDecoder struct{}

func NewInternalImageDecoder() *InternalImageDecoder {
	return &InternalImageDecoder{}
}

func (i *InternalImageDecoder) Decode(reader io.Reader) (image.Image, error) {
	return jpeg.Decode(reader) //nolint: wrapcheck
}
