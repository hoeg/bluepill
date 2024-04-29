package adaptors

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/g4s8/go-lifecycle/pkg/adaptors"
	"github.com/g4s8/go-lifecycle/pkg/types"
	"github.com/pkg/errors"
)

// HTTPService is an adapter for http.Server to implement lifecycle hooks.
type HTTPSService struct {
	srv      *http.Server
	certFile string
	keyFile  string
}

// NewHTTPService creates new HTTPService adaptor.
func NewHTTPSService(srv *http.Server, certFile string, keyFile string) *HTTPSService {
	return &HTTPSService{
		srv:      srv,
		certFile: certFile,
		keyFile:  keyFile,
	}
}

// RegisterLifecycle registers this service in lifecycle manager.
func (s *HTTPSService) RegisterLifecycle(name string, lf adaptors.LifecycleRegistry) {
	lf.RegisterService(types.ServiceConfig{
		Name:         name,
		StartupHook:  s.Start,
		ShutdownHook: s.Stop,
		RestartPolicy: types.ServiceRestartPolicy{
			RestartOnFailure: true,
			RestartCount:     3,
			RestartDelay:     time.Millisecond * 200,
		},
	})
}

func (s *HTTPSService) Start(ctx context.Context, errCh chan<- error) error {
	addr := s.srv.Addr
	if addr == "" {
		addr = ":https"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen tcp: %w", err)
	}
	go func() {
		if err := s.srv.ServeTLS(ln, s.certFile, s.keyFile); err != nil {
			log.Printf("serve error: %v", err)
			if err != http.ErrServerClosed {
				errCh <- errors.Wrap(err, "serve")
				return
			}
		}
	}()
	return nil
}

func (s *HTTPSService) Stop(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shutdown")
	}
	return nil
}
