package app

import (
	"log"

	"github.com/g4s8/go-lifecycle/pkg/lifecycle"
	"github.com/hoeg/bluepill/internal/morpheus"
)

func Start() {
	conf, err := LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
		return
	}
	api := morpheus.NewAPI(conf.enforcementConfig)
	svc := NewHTTPService(api, conf.HTTPConfig)

	lf := lifecycle.New(lifecycle.DefaultConfig)
	svc.RegisterLifecycle("web", lf)

	lf.Start()
	sig := lifecycle.NewSignalHandler(lf, nil)
	sig.Start(lifecycle.DefaultShutdownConfig)
	if err := sig.Wait(); err != nil {
		log.Fatalf("shutdown error: %v", err)
		return
	}
	log.Print("shutdown complete")
}
