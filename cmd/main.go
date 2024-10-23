package main

import (
	"fmt"
	"log/slog"

	"github.com/emiliocc5/CheckBoxDetector/config"
	"github.com/emiliocc5/CheckBoxDetector/internal/adapters"
	"github.com/emiliocc5/CheckBoxDetector/internal/core"
	"github.com/emiliocc5/CheckBoxDetector/internal/ports"
)

func main() {
	log := slog.Default()
	detector := initialize()

	resp, errDetection := detector.Detect()
	if errDetection != nil {
		log.Error("error on checkbox detection")
	}

	log.Info(fmt.Sprintf("Total of detected checkboxes is: %d", resp))
}

func initialize() ports.CheckBoxDetector {
	return core.NewAutomaticDetector(
		getConfig(),
		adapters.NewInternalImageGetter(),
		adapters.NewInternalImageDecoder(),
		adapters.NewInternalImageGrayer(),
		adapters.NewInternalImageBinarizer())
}

func getConfig() config.ServiceConfiguration {
	conf := config.NewConfigService()
	conf.LoadFromEnvs()

	return *conf
}
