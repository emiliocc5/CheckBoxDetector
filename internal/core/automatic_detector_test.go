package core

import (
	"errors"
	"github.com/emiliocc5/CheckBoxDetector/internal/mocks"
	"github.com/golang/mock/gomock"
	"testing"

	"github.com/emiliocc5/CheckBoxDetector/config"
	"github.com/emiliocc5/CheckBoxDetector/internal/adapters"
	"github.com/stretchr/testify/assert"
)

func Test_AutomaticDetectImage_SuccessDetection(t *testing.T) {
	t.Parallel()
	cfg := config.NewConfigService()
	cfg.File.Path = "../../resources/form.jpg"

	detector := NewAutomaticDetector(
		*cfg,
		adapters.NewInternalImageGetter(),
		adapters.NewInternalImageDecoder(),
		adapters.NewInternalImageGrayer(),
		adapters.NewInternalImagePixelHandler(),
		adapters.NewInternalImageBinarizer())

	resp, err := detector.Detect()
	assert.NoError(t, err)
	assert.Equal(t, 15, resp)
}

func Test_AutomaticDetectImage_ErrorGettingImage(t *testing.T) {
	t.Parallel()
	cfg := config.NewConfigService()

	ctrl := gomock.NewController(t)
	mockInternalImageGetter := mocks.NewMockImageGetter(ctrl)
	mockInternalImageGetter.EXPECT().GetImage(gomock.Any()).Return(nil,
		errors.New("unexpected error getting image"))

	detector := NewAutomaticDetector(
		*cfg,
		mockInternalImageGetter,
		adapters.NewInternalImageDecoder(),
		adapters.NewInternalImageGrayer(),
		adapters.NewInternalImagePixelHandler(),
		adapters.NewInternalImageBinarizer())

	resp, err := detector.Detect()
	assert.Equal(t, errors.New("error getting image"), err)
	assert.Equal(t, 0, resp)
}

func Test_AutomaticDetectImage_ErrorDecodingImage(t *testing.T) {
	t.Parallel()
	cfg := config.NewConfigService()
	cfg.File.Path = "../../resources/form.jpg"

	ctrl := gomock.NewController(t)
	mockInternalImageDecoder := mocks.NewMockImageDecoder(ctrl)
	mockInternalImageDecoder.EXPECT().Decode(gomock.Any()).Return(nil,
		errors.New("unexpected error decoding image"))

	detector := NewAutomaticDetector(
		*cfg,
		adapters.NewInternalImageGetter(),
		mockInternalImageDecoder,
		adapters.NewInternalImageGrayer(),
		adapters.NewInternalImagePixelHandler(),
		adapters.NewInternalImageBinarizer())

	resp, err := detector.Detect()
	assert.Equal(t, errors.New("error decoding image"), err)
	assert.Equal(t, 0, resp)
}
