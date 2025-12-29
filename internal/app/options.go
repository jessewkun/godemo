package app

import (
	"fmt"

	"github.com/jessewkun/gocommon/alarm"
	"github.com/jessewkun/gocommon/config"
	"github.com/jessewkun/gocommon/logger"
)

// Options holds the application's configuration.
type Options struct {
	ConfigFile string
	BaseConfig *config.BaseConfig
}

// NewOptions creates a new Options instance by loading the application configuration.
// It expects the configFile path to be provided.
func NewOptions(configFile string) (*Options, error) {
	o := &Options{
		ConfigFile: configFile,
	}

	if err := o.loadConfig(); err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %w", configFile, err)
	}

	return o, nil
}

// crossModuleInjector implements cross-module dependency injection,
// preventing circular dependencies at the lower levels.
func crossModuleInjector() error {
	// Inject the alarm implementation into the logger.
	var alarmSender alarm.Sender
	logger.RegisterAlerter(&alarmSender)
	fmt.Println("Cross-module injection for logger and alarm completed successfully.")
	return nil
}

// loadConfig loads the configuration from the configuration file.
func (o *Options) loadConfig() error {
	config.RegisterInjector(crossModuleInjector)

	baseConfig, err := config.Init(o.ConfigFile)
	if err != nil {
		return fmt.Errorf("load config file %s error: %w", o.ConfigFile, err)
	}
	o.BaseConfig = baseConfig
	return nil
}
