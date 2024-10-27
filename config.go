package kqueuey

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var (
	TestFlag    uint
	SearchPaths = []string{SearchPath0, SearchPath1, SearchPath2, SearchPath3, SearchPath4}

	errConfigFileNotLocated = fmt.Errorf("config file not found, valid paths: %v", SearchPaths)
)

const (
	SearchPath0 = "/config/"
	SearchPath1 = "/etc/config/"
	SearchPath2 = "/var/lib/config/"
	SearchPath3 = "/data/"
	SearchPath4 = "config.yaml"

	ConfigType = "yaml"
)

type Configuration struct {
	BadgerOpts BadgerDB `mapstructure:"storage"`
	RaftOpts   Raft     `mapstructure:"raft"`
}

type BadgerDB struct {
	NumCompactors   int    `mapstructure:"num_compactors"`
	CompressionType string `mapstructure:"compression_type"`
	SyncWrites      bool   `mapstructure:"sync_writes"`
}

type Raft struct {
	ClusterID string `mapstructure:"cluster_id"`
	Nodes     []Node `mapstructure:"nodes"`
}

type Node struct {
	Id         string `mapstructure:"id"`
	BindAddr   string `mapstructure:"bind_addr"`
	StorageDir string `mapstructure:"storage_dir"`
}

func LoadConfiguration() (*Configuration, error) {
	v := getViper()

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	if len(v.ConfigFileUsed()) == 0 {
		return nil, errConfigFileNotLocated
	}

	var config Configuration
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	if err = config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func getViper() *viper.Viper {
	v := viper.New()
	setDefaults(v)
	v.SetConfigType(ConfigType)

	if TestFlag == 1 {
		pwd, err := os.Getwd()
		if err == nil {
			v.AddConfigPath(pwd + "/integration/")
		}
		v.AddConfigPath("config.yaml")
	} else {
		for _, path := range SearchPaths {
			v.AddConfigPath(path)
		}
	}

	return v
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("storage.num_compactors", 4)
	v.SetDefault("storage.compression_type", CompressionSnappy)
}

// Validate ensures that the config file is set up correctly, to allow proper start up of queue server.
func (c *Configuration) Validate() error {
	if c.BadgerOpts.NumCompactors < 4 {
		c.BadgerOpts.NumCompactors = 4
	}

	if c.RaftOpts.ClusterID == "" {
		return errors.New("cluster id was not found")
	}

	duplicates := make(map[string]uint8)
	for _, node := range c.RaftOpts.Nodes {
		splitAddr := strings.Split(node.BindAddr, ":")
		duplicates[node.StorageDir] += 1
		duplicates[node.Id] += 1

		if duplicates[node.StorageDir] >= 2 {
			return errors.New("storage path must be unique per node to avoid file locking conflicts")
		}

		if duplicates[node.Id] >= 2 {
			return errors.New("node id must be unique per node to conform to rafts consensus protocol")
		}

		if node.Id == "" {
			return errors.New("node id not found")
		}

		if node.StorageDir == "" {
			return errors.New("node id not found")
		}

		if len(splitAddr) != 2 || len(splitAddr[1]) < 4 {
			return errors.New("unknown raft address")
		}

		duplicates[splitAddr[1]] += 1
		if duplicates[splitAddr[1]] == 2 {
			return errors.New("node port most be unique to allow proper communication between nodes")
		}
	}

	return nil
}
