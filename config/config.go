package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Verbose bool
	LogPath string
	JSON    bool
	DBURL   string // e.g. file:replicator.db?cache=shared&_busy_timeout=5000
}

type fileConfig struct {
	Log struct {
		Path    string `toml:"path"`
		JSON    bool   `toml:"json"`
		Verbose bool   `toml:"verbose"`
	} `toml:"log"`
	Database struct {
		URL string `toml:"url"`
	} `toml:"database"`
}

const (
	defaultConfigPath = "config.toml"
	defaultDBURL      = "file:replicator.db?cache=shared&_busy_timeout=5000"
)

func LoadConfig() *Config {
	var cfgPath string
	var verbose bool

	flag.StringVar(&cfgPath, "config", "", "Path to config TOML")
	flag.BoolVar(&verbose, "v", false, "Enable verbose logging")
	flag.Parse()

	if cfgPath == "" {
		cfgPath = defaultConfigPath
	}
	if _, err := os.Stat(cfgPath); err != nil {
		panic(fmt.Sprintf("config file not found: %s", cfgPath))
	}

	var fc fileConfig
	md, err := toml.DecodeFile(cfgPath, &fc)
	if err != nil {
		panic(fmt.Errorf("failed to parse config: %w", err))
	}
	if undec := md.Undecoded(); len(undec) > 0 {
		keys := make([]string, 0, len(undec))
		for _, k := range undec {
			keys = append(keys, k.String())
		}
		panic("unknown config keys: " + strings.Join(keys, ", "))
	}

	c := &Config{
		Verbose: verbose,
		LogPath: "",
		JSON:    false,
		DBURL:   defaultDBURL,
	}
	c.Verbose = fc.Log.Verbose
	if fc.Log.Path != "" {
		c.LogPath = fc.Log.Path
	}
	c.JSON = fc.Log.JSON
	if fc.Database.URL != "" {
		c.DBURL = fc.Database.URL
	}

	return c
}
