package kqueuey

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var (
	TestFlag    uint
	SearchPaths = []string{SearchPath0, SearchPath1, SearchPath2, SearchPath3}

	errConfigFileNotLocated   = fmt.Errorf("config file not found, valid paths: %v", SearchPaths)
	errDuplicateStoragePath   = errors.New("storage path must be unique per node to avoid file locking conflicts")
	errDuplicateNodeId        = errors.New("node id must be unique per node to conform to rafts consensus protocol")
	errDuplicateNodePort      = errors.New("node port most be unique to allow proper communication between nodes")
	errNodeIdNotFound         = errors.New("node id not found")
	errNodeStorageDirNotFound = errors.New("node storage directory not found")
	errUnknownNodeAddress     = errors.New("unknown node address")
	errRaftClusterIdNotFound  = errors.New("cluster id was not found")
)

const (
	SearchPath0 = "/config/"
	SearchPath1 = "/etc/config/"
	SearchPath2 = "/var/lib/config/"
	SearchPath3 = "/data/"

	ConfigType     = "yaml"
	ConfigFileName = "kqueuey-config"
)

type Configuration struct {
	BadgerOpts BadgerDBConfig `mapstructure:"storage"`
	RaftOpts   RaftConfig     `mapstructure:"raft"`
}

type BadgerDBConfig struct {
	NumCompactors   int    `mapstructure:"num_compactors"`
	CompressionType string `mapstructure:"compression_type"`
	SyncWrites      bool   `mapstructure:"sync_writes"`
}

type RaftConfig struct {
	ClusterId string       `mapstructure:"cluster_id"`
	Nodes     []NodeConfig `mapstructure:"nodes"`
}

type NodeConfig struct {
	Id         string `mapstructure:"id"`
	BindAddr   string `mapstructure:"bind_addr"`
	StorageDir string `mapstructure:"storage_dir"`
}

func LoadConfiguration() (*Configuration, error) {
	v := initializeViper()

	if err := v.ReadInConfig(); err != nil {
		return nil, errConfigFileNotLocated
	}

	var config Configuration
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func initializeViper() *viper.Viper {
	v := viper.New()
	setDefaults(v)
	v.SetConfigType(ConfigType)
	v.SetConfigName(ConfigFileName)

	for _, path := range SearchPaths {
		v.AddConfigPath(path)
	}

	return v
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("storage.num_compactors", 4)
	v.SetDefault("storage.compression_type", CompressionSnappy)
}

// Validate ensures that the config file is set up correctly, to allow proper start up of queue server.
func (c *Configuration) validate() error {
	if c.BadgerOpts.NumCompactors < 4 {
		c.BadgerOpts.NumCompactors = 4
	}

	if err := c.RaftOpts.validateRaftConfigOptions(); err != nil {
		return err
	}

	return nil
}

func (r *RaftConfig) validateRaftConfigOptions() error {
	if r.ClusterId == "" {
		return errRaftClusterIdNotFound
	}

	duplicates := make(map[string]uint8)
	for _, node := range r.Nodes {
		if node.Id == "" {
			return errNodeIdNotFound
		}

		if node.StorageDir == "" {
			return errNodeStorageDirNotFound
		}

		duplicates[node.StorageDir] += 1
		if duplicates[node.StorageDir] >= 2 {
			return errDuplicateStoragePath
		}

		duplicates[node.Id] += 1
		if duplicates[node.Id] >= 2 {
			return errDuplicateNodeId
		}

		splitAddr := strings.Split(node.BindAddr, ":")
		port := splitAddr[1]
		if len(splitAddr) != 2 || len(port) < 4 {
			return errUnknownNodeAddress
		}

		duplicates[port] += 1
		if duplicates[port] == 2 {
			return errDuplicateNodePort
		}
	}

	return nil
}
