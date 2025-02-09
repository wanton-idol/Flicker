package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/SuperMatch/config"
)

type Server struct {
	s       *http.Server
	closeFn func()
}

func Start(router http.Handler, conf *config.HTTP, closeFn func()) (Server, error) {

	s := Server{
		closeFn: closeFn,
	}

	if conf.PORT == "" {
		return s, errors.New("ENV PORT is not defined.must be a valid TCP PORT")
	}
	s.s = &http.Server{
		Addr:         "localhost:" + conf.PORT,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		err := s.s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.closeFn()
			log.Fatalf("http error :%v\n", err)
		}
	}()

	return s, nil
}

func (s *Server) Stop() error {
	err := s.s.Shutdown(context.Background())
	s.closeFn()
	return err
}

func StartServer(router http.Handler, htttpConf config.HTTP, wg *sync.WaitGroup) Server {
	wg.Add(1)

	s, err := Start(router, &htttpConf, func() {
		wg.Done()
	})

	if err != nil {
		log.Fatalf("http error :%v\n", err)
	}
	return s
}
