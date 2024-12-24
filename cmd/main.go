package main

import (
	"fmt"
	"github.com/emiliocc5/CheckBoxDetector/config"
	"github.com/emiliocc5/CheckBoxDetector/internal/adapters"
	"github.com/emiliocc5/CheckBoxDetector/internal/core"
	"github.com/emiliocc5/CheckBoxDetector/internal/ports"
	"log/slog"
	"sync"
)

func main() {
	log := slog.Default()
	detector := initialize()
	ch := make(chan map[int]int)
	errCh := make(chan error)
	var wg sync.WaitGroup

	for i := 1; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, errDetection := detector.Detect(i)
			if errDetection != nil {
				log.Error("error on checkbox detection")
				errCh <- errDetection
			}
			ch <- map[int]int{i: resp}
		}()
		//log.Info(fmt.Sprintf("Total of detected checkboxes is: %d", resp))
	}

	go func() {
		wg.Wait()
		close(ch)
		close(errCh)
	}()

	for resp := range ch {
		log.Info(fmt.Sprintf("Total of detected checkboxes is: %d", resp))
	}

	for errors := range errCh {
		fmt.Println(errors)
	}
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
