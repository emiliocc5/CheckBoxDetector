package adapters

import (
	"io"
	"log/slog"
	"os"
)

type InternalImageGetter struct{}

func NewInternalImageGetter() *InternalImageGetter {
	return &InternalImageGetter{}
}

func (i *InternalImageGetter) GetImage(path string) (io.Reader, error) {
	log := slog.Default()
	file, err := os.Open(path)
	if err != nil {
		log.Error("error loading image", slog.String("error", err.Error()))

		return nil, err //nolint: wrapcheck
	}

	return file, nil
}
