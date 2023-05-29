package server

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/rs/zerolog/log"
	"github.com/rytsh/liz/utils/shutdown"
	"github.com/rytsh/yap/internal/tui"
)

var ShutdownTimeout = 5 * time.Second

type Config struct {
	Host string
	Port int

	Screen tui.Screen
}

func Serve(wg *sync.WaitGroup, cfg Config) error {
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			screenMiddleware(cfg.Screen),
			lm.Middleware(),
		),
	)
	if err != nil {
		return fmt.Errorf("could not create server: %w", err)
	}

	log.Info().Msgf("starting yap server on %s", s.Addr)

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error().Err(err).Msg("could not start server")
		}
	}()

	shutdown.Global.Add("yap server", func() error {
		return shutdownServer(s)
	})

	return nil
}

func shutdownServer(s *ssh.Server) error {
	log.Info().Msg("Stopping SSH server")

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		return fmt.Errorf("could not stop server: %w", err)
	}

	return nil
}
