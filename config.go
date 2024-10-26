package kqueuey

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	TestFlag    uint8
	SearchPaths = []string{SearchPath0, SearchPath1, SearchPath2, SearchPath3}
)

const (
	SearchPath0 = "/config/"
	SearchPath1 = "/etc/config/"
	SearchPath2 = "/var/lib/config/"
	SearchPath3 = "/data/"
)

type Configuration struct {
	BadgerOpts BadgerDB `mapstructure:"badgerDB"`
	RaftOpts   Raft     `mapstructure:"raft"`
}

type BadgerDB struct {
	StorePath       []string `mapstructure:"kv-store-path"`
	NumCompactors   int      `mapstructure:"num-compactors"`
	CompressionType string   `mapstructure:"compression-type"`
	SyncWrites      bool     `mapstructure:"sync-writes"`
}

type Raft struct {
	Nodes []map[string]any `mapstructure:"nodes"`
}

func LoadConfiguration() (*Configuration, error) {
	v := getViper()

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Configuration
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func getViper() *viper.Viper {
	v := viper.New()
	setDefault(v)
	v.SetConfigType("yaml")
	v.SetConfigFile("kqueuey-config.yaml")

	// Testing path.
	if TestFlag == 1 {
		homeDir, _ := os.UserHomeDir()
		v.AddConfigPath(fmt.Sprintf("%s", homeDir))
		v.AddConfigPath("/integration/")
	} else {
		// Currently acceptable paths to obtain kqueuey config.
		for _, path := range SearchPaths {
			v.AddConfigPath(path)
		}
	}

	watchForConfigUpdates(v)

	return v
}

// watch enables hot reloading to accommodate changes in the config.
// Such as scaling up nodes for raft or updating existing values.
// This ensures that raft can dynamically recognize new nodes joining the cluster while the server remains active.
func watchForConfigUpdates(c *viper.Viper) {
	// Send a ReadinConfig notice to allow viper to re-load the configuration file from disk when a change occurred.
	c.WatchConfig()
	c.OnConfigChange(func(e fsnotify.Event) {
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Println("issue")
		}
	})
}

func setDefault(v *viper.Viper) {
	v.SetDefault("badgerDB.num-compactors", 4)
	v.SetDefault("badgerDB.compression-type", Snappy)
}
