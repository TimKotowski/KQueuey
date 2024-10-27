package kqueuey

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type StorageConfigTest struct {
	NumCompactors   int    `yaml:"num_compactors"`
	CompressionType string `yaml:"compression_type"`
	SyncWrites      bool   `yaml:"sync_writes"`
}

type RaftNodeTest struct {
	ID         string `yaml:"id"`
	BindAddr   string `yaml:"bind_addr"`
	StorageDir string `yaml:"storage_dir"`
}

type RaftConfigTest struct {
	ClusterID string         `yaml:"cluster_id"`
	Nodes     []RaftNodeTest `yaml:"nodes"`
}

type ConfigTest struct {
	Storage StorageConfigTest `yaml:"storage"`
	Raft    RaftConfigTest    `yaml:"raft"`
}

func TestViperConfig(t *testing.T) {
	t.Run("read and parse kqueuey config with all correct values ", func(t *testing.T) {
		defer configFileCleanUp()

		TestFlag = 1
		expectedStorageConfig := StorageConfigTest{NumCompactors: 1, CompressionType: "snappy", SyncWrites: false}
		expectedRaftConfig := RaftConfigTest{
			ClusterID: "1000",
			Nodes: []RaftNodeTest{
				{
					ID:         "1",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path0/dir",
				},

				{
					ID:         "2",
					BindAddr:   "127.0.0.1:1211",
					StorageDir: "path1/dir",
				},

				{
					ID:         "3",
					BindAddr:   "127.0.0.1:1311",
					StorageDir: "path2/dir",
				},
			},
		}
		setUpTestConfigFile(t, expectedStorageConfig, expectedRaftConfig)
		config, err := LoadConfiguration()
		assert.NoError(t, err)
		assert.EqualValues(t, expectedStorageConfig.SyncWrites, config.BadgerOpts.SyncWrites)
		assert.Equal(t, config.BadgerOpts.NumCompactors, 4)
		assert.Equal(t, expectedStorageConfig.CompressionType, config.BadgerOpts.CompressionType)

		expectedNodeOne := expectedRaftConfig.Nodes[0]
		expectedNodeTwo := expectedRaftConfig.Nodes[1]
		expectedNodeThree := expectedRaftConfig.Nodes[2]
		actualNodeOne := config.RaftOpts.Nodes[0]
		actualNodeTwo := config.RaftOpts.Nodes[1]
		actualNodeThree := config.RaftOpts.Nodes[2]
		assert.Equal(t, expectedRaftConfig.ClusterID, config.RaftOpts.ClusterID)
		assert.Equal(t, len(expectedRaftConfig.Nodes), len(config.RaftOpts.Nodes))
		assert.Equal(t, expectedRaftConfig.ClusterID, config.RaftOpts.ClusterID)

		assert.Equal(t, expectedNodeOne.ID, actualNodeOne.Id)
		assert.Equal(t, expectedNodeTwo.ID, actualNodeTwo.Id)
		assert.Equal(t, expectedNodeThree.ID, actualNodeThree.Id)

		assert.Equal(t, expectedNodeOne.BindAddr, actualNodeOne.BindAddr)
		assert.Equal(t, expectedNodeTwo.BindAddr, actualNodeTwo.BindAddr)
		assert.Equal(t, expectedNodeThree.BindAddr, actualNodeThree.BindAddr)

		assert.Equal(t, expectedNodeOne.StorageDir, actualNodeOne.StorageDir)
		assert.Equal(t, expectedNodeTwo.StorageDir, actualNodeTwo.StorageDir)
		assert.Equal(t, expectedNodeThree.StorageDir, actualNodeThree.StorageDir)
	})
	t.Run("kqueuey configuration file validation failed due to duplicate node id", func(t *testing.T) {
		defer configFileCleanUp()
		TestFlag = 1

		expectedStorageConfig := StorageConfigTest{NumCompactors: 1, CompressionType: "snappy", SyncWrites: false}
		expectedRaftConfig := RaftConfigTest{
			ClusterID: "1000",
			Nodes: []RaftNodeTest{
				{
					ID:         "1",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path0/dir",
				},

				{
					ID:         "1",
					BindAddr:   "127.0.0.1:1211",
					StorageDir: "path1/dir",
				},

				{
					ID:         "3",
					BindAddr:   "127.0.0.1:1311",
					StorageDir: "path2/dir",
				},
			},
		}
		setUpTestConfigFile(t, expectedStorageConfig, expectedRaftConfig)
		config, err := LoadConfiguration()
		assert.Error(t, err)
		assert.ErrorIs(t, err, errDuplicateNodeID)
		assert.Nil(t, config)
	})

	t.Run("kqueuey configuration file validation failed due to duplicate node port", func(t *testing.T) {
		defer configFileCleanUp()
		TestFlag = 1

		expectedStorageConfig := StorageConfigTest{NumCompactors: 1, CompressionType: "snappy", SyncWrites: false}
		expectedRaftConfig := RaftConfigTest{
			ClusterID: "1000",
			Nodes: []RaftNodeTest{
				{
					ID:         "1",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path0/dir",
				},

				{
					ID:         "2",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path1/dir",
				},

				{
					ID:         "3",
					BindAddr:   "127.0.0.1:1311",
					StorageDir: "path2/dir",
				},
			},
		}
		setUpTestConfigFile(t, expectedStorageConfig, expectedRaftConfig)
		config, err := LoadConfiguration()
		assert.Error(t, err)
		assert.ErrorIs(t, err, errDuplicateNodePort)
		assert.Nil(t, config)
	})

	t.Run("kqueuey configuration file validation failed due to duplicate node port", func(t *testing.T) {
		defer configFileCleanUp()
		TestFlag = 1

		expectedStorageConfig := StorageConfigTest{NumCompactors: 1, CompressionType: "snappy", SyncWrites: false}
		expectedRaftConfig := RaftConfigTest{
			ClusterID: "1000",
			Nodes: []RaftNodeTest{
				{
					ID:         "1",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path0/dir",
				},

				{
					ID:         "1",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path0/dir",
				},

				{
					ID:         "3",
					BindAddr:   "127.0.0.1:1311",
					StorageDir: "path2/dir",
				},
			},
		}
		setUpTestConfigFile(t, expectedStorageConfig, expectedRaftConfig)
		config, err := LoadConfiguration()
		assert.Error(t, err)
		assert.ErrorIs(t, err, errDuplicateStoragePath)
		assert.Nil(t, config)
	})

	t.Run("kqueuey configuration file validation failed due to duplicate cluster id", func(t *testing.T) {
		defer configFileCleanUp()
		TestFlag = 1

		expectedStorageConfig := StorageConfigTest{NumCompactors: 1, CompressionType: "snappy", SyncWrites: false}
		expectedRaftConfig := RaftConfigTest{
			ClusterID: "",
			Nodes: []RaftNodeTest{
				{
					ID:         "1",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path0/dir",
				},

				{
					ID:         "1",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path0/dir",
				},

				{
					ID:         "3",
					BindAddr:   "127.0.0.1:1311",
					StorageDir: "path2/dir",
				},
			},
		}
		setUpTestConfigFile(t, expectedStorageConfig, expectedRaftConfig)
		config, err := LoadConfiguration()
		assert.Error(t, err)
		assert.ErrorIs(t, err, errRaftClusterIdNotFound)
		assert.Nil(t, config)
	})

	t.Run("kqueuey configuration file validation failed due to unkown node address", func(t *testing.T) {
		defer configFileCleanUp()
		TestFlag = 1

		expectedStorageConfig := StorageConfigTest{NumCompactors: 1, CompressionType: "snappy", SyncWrites: false}
		expectedRaftConfig := RaftConfigTest{
			ClusterID: "1000",
			Nodes: []RaftNodeTest{
				{
					ID:         "1",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path0/dir",
				},

				{
					ID:         "2",
					BindAddr:   "127.0.0.1:",
					StorageDir: "path1/dir",
				},

				{
					ID:         "3",
					BindAddr:   "127.0.0.1:1",
					StorageDir: "path2/dir",
				},
			},
		}
		setUpTestConfigFile(t, expectedStorageConfig, expectedRaftConfig)
		config, err := LoadConfiguration()
		assert.Error(t, err)
		assert.ErrorIs(t, err, errUnknownNodeAddress)
		assert.Nil(t, config)
	})

	t.Run("kqueuey configuration file validation failed due to missing node id", func(t *testing.T) {
		defer configFileCleanUp()
		TestFlag = 1

		expectedStorageConfig := StorageConfigTest{NumCompactors: 1, CompressionType: "snappy", SyncWrites: false}
		expectedRaftConfig := RaftConfigTest{
			ClusterID: "1000",
			Nodes: []RaftNodeTest{
				{
					ID:         "",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "path0/dir",
				},

				{
					ID:         "1",
					BindAddr:   "127.0.0.1:1211",
					StorageDir: "path0/dir",
				},

				{
					ID:         "3",
					BindAddr:   "127.0.0.1:1311",
					StorageDir: "path2/dir",
				},
			},
		}
		setUpTestConfigFile(t, expectedStorageConfig, expectedRaftConfig)
		config, err := LoadConfiguration()
		assert.Error(t, err)
		assert.ErrorIs(t, err, errNodeIdNotFound)
		assert.Nil(t, config)
	})

	t.Run("kqueuey configuration file validation failed due missing storage directory for node", func(t *testing.T) {
		defer configFileCleanUp()
		TestFlag = 1

		expectedStorageConfig := StorageConfigTest{NumCompactors: 1, CompressionType: "snappy", SyncWrites: false}
		expectedRaftConfig := RaftConfigTest{
			ClusterID: "1000",
			Nodes: []RaftNodeTest{
				{
					ID:         "1",
					BindAddr:   "127.0.0.1:1111",
					StorageDir: "",
				},

				{
					ID:         "1",
					BindAddr:   "127.0.0.1:1211",
					StorageDir: "path0/dir",
				},

				{
					ID:         "3",
					BindAddr:   "127.0.0.1:1311",
					StorageDir: "path2/dir",
				},
			},
		}
		setUpTestConfigFile(t, expectedStorageConfig, expectedRaftConfig)
		config, err := LoadConfiguration()
		assert.Error(t, err)
		assert.ErrorIs(t, err, errNodeStorageDirNotFound)
		assert.Nil(t, config)
	})
}

func setUpTestConfigFile(t *testing.T, storageConfig StorageConfigTest, raftConfig RaftConfigTest) {
	c := ConfigTest{Storage: storageConfig, Raft: raftConfig}
	o, err := yaml.Marshal(&c)
	assert.NoError(t, err)
	err = os.WriteFile("kqueuey-config.yaml", o, 0644)
	assert.NoError(t, err)
}

func configFileCleanUp() {
	_ = os.Remove("kqueuey-config.yaml")
}
