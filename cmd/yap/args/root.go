package args

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/rytsh/liz/utils/shutdown"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/worldline-go/igconfig"
	"github.com/worldline-go/igconfig/loader"
	"github.com/worldline-go/logz"

	"github.com/rytsh/yap/internal/config"
	"github.com/rytsh/yap/internal/server"
)

var ErrShutdown = errors.New("shutting down signal received")

var rootCmd = &cobra.Command{
	Use:   "yap",
	Short: "run codes",
	Long:  fmt.Sprintf("%s\nyap run codes", config.GetBanner()),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := logz.SetLogLevel(config.Application.LogLevel); err != nil {
			return err //nolint:wrapcheck // no need
		}

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// load configuration
		if err := loadConfig(cmd.Context(), cmd.Flags().Visit); err != nil {
			return err
		}

		if err := runRoot(cmd.Context()); err != nil && !errors.Is(err, ErrShutdown) {
			return err
		}

		return nil
	},
}

func Execute(ctx context.Context) error {
	setFlags()

	rootCmd.Version = config.BuildVars.Version
	rootCmd.Long = fmt.Sprintf(
		"%s\nversion:[%s] commit:[%s] buildDate:[%s]",
		rootCmd.Long, config.BuildVars.Version, config.BuildVars.Commit, config.BuildVars.Date,
	)

	return rootCmd.ExecuteContext(ctx) //nolint:wrapcheck // no need
}

func setFlags() {
	rootCmd.PersistentFlags().StringVarP(&config.Application.LogLevel, "log-level", "l", config.Application.LogLevel, "log level")
	rootCmd.PersistentFlags().StringVarP(&config.Application.Server.Host, "host", "H", config.Application.Server.Host, "host")
	rootCmd.PersistentFlags().IntVarP(&config.Application.Server.Port, "port", "P", config.Application.Server.Port, "port")
}

// override function hold first values of definitions.
// Use with pflag visit function.
func override(ow map[string]func()) {
	ow["log-level"] = func(v string) func() { return func() { config.Application.LogLevel = v } }(config.Application.LogLevel)
	ow["host"] = func(v string) func() { return func() { config.Application.Server.Host = v } }(config.Application.Server.Host)
	ow["port"] = func(v int) func() { return func() { config.Application.Server.Port = v } }(config.Application.Server.Port)
}

func loadConfig(ctx context.Context, visit func(fn func(*pflag.Flag))) error {
	overrideValues := make(map[string]func())
	override(overrideValues)

	logConfig := log.With().Str("component", "config").Logger()
	ctxConfig := logConfig.WithContext(ctx)

	loaders := []loader.Loader{}

	envLoader := &loader.Env{}

	if err := igconfig.LoadWithLoadersWithContext(ctxConfig, "", &config.LoadConfig, envLoader); err != nil {
		return fmt.Errorf("unable to load prefix settings: %v", err)
	}

	log.Info().Msgf("config loading from %+v", config.LoadConfig)

	loader.ConsulConfigPathPrefix = config.LoadConfig.Prefix.Consul
	loader.VaultSecretBasePath = config.LoadConfig.Prefix.Vault
	loader.VaultSecretAdditionalPaths = nil

	if config.LoadConfig.ConfigSet.Consul {
		loaders = append(loaders, &loader.Consul{})
	}

	if config.LoadConfig.ConfigSet.Vault && config.LoadConfig.Prefix.Vault != "" {
		loaders = append(loaders, &loader.Vault{})
	}

	if config.LoadConfig.ConfigSet.File {
		loaders = append(loaders, &loader.File{})
	}

	loaders = append(loaders, envLoader)

	if err := igconfig.LoadWithLoadersWithContext(ctxConfig, config.LoadConfig.AppName, &config.Application, loaders...); err != nil {
		return fmt.Errorf("unable to load configuration settings: %v", err)
	}

	// override used cmd values
	visit(func(f *pflag.Flag) {
		if v, ok := overrideValues[f.Name]; ok {
			v()
		}
	})

	// set log again to get changes
	if err := logz.SetLogLevel(config.Application.LogLevel); err != nil {
		return err //nolint:wrapcheck // no need
	}

	// print loaded object
	log.Debug().Object("config", igconfig.Printer{Value: config.Application}).Msg("loaded config")

	return nil
}

func runRoot(ctxParent context.Context) (err error) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	ctx, ctxCancel := context.WithCancel(ctxParent)
	defer ctxCancel()

	wg.Add(1)

	go func() {
		defer wg.Done()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-sig:
			log.Warn().Msg("received shutdown signal")
			ctxCancel()

			if err != nil {
				err = ErrShutdown
			}
		case <-ctx.Done():
			log.Warn().Msg("yap closing")
		}

		shutdown.Global.Run()
	}()

	// application codes
	if err := server.Serve(wg, server.Config{
		Host:   config.Application.Server.Host,
		Port:   config.Application.Server.Port,
		Screen: config.Application.Screen,
	}); err != nil {
		return err
	}

	wg.Wait()

	return nil
}
