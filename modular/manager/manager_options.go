package manager

import (
	"github.com/bnb-chain/greenfield-storage-provider/base/gfspapp"
	"github.com/bnb-chain/greenfield-storage-provider/base/gfspconfig"
	coremodule "github.com/bnb-chain/greenfield-storage-provider/core/module"
)

const (
	// DefaultGlobalMaxUploadingNumber defines the default max uploading object number
	// in SP, include: uploading object to primary, replicate object to secondaries,
	// and sealing object on greenfield.
	DefaultGlobalMaxUploadingNumber int = 4096
	// DefaultGlobalUploadObjectParallel defines the default max parallel uploading
	// objects to primary in SP system.
	DefaultGlobalUploadObjectParallel int = 1024
	// DefaultGlobalReplicatePieceParallel defines the default max parallel replicating
	// objects to primary in SP system.
	DefaultGlobalReplicatePieceParallel int = 1024
	// DefaultGlobalSealObjectParallel defines the default max parallel sealing objects
	// on greenfield in SP system.
	DefaultGlobalSealObjectParallel int = 1024
	// DefaultGlobalReceiveObjectParallel defines the default max parallel confirming
	// receive pieces on greenfield in SP system.
	DefaultGlobalReceiveObjectParallel int = 4096
	// DefaultGlobalGCObjectParallel defines the default max parallel gc objects in SP
	// system.
	DefaultGlobalGCObjectParallel int = 4
	// DefaultGlobalGCZombieParallel defines the default max parallel gc zonbie pieces
	// in SP system.
	DefaultGlobalGCZombieParallel int = 1
	// DefaultGlobalGCMetaParallel defines the default max parallel gc meta db in SP
	// system.
	DefaultGlobalGCMetaParallel int = 1
	// DefaultGlobalDownloadObjectTaskCacheSize defines the default max cache the download
	// object tasks in manager.
	DefaultGlobalDownloadObjectTaskCacheSize int = 4096
	// DefaultGlobalChallengePieceTaskCacheSize defines the default max cache the challenge
	// piece tasks in manager.
	DefaultGlobalChallengePieceTaskCacheSize int = 4096
	// DefaultGlobalBatchGcObjectTimeInterval defines the default interval for generating
	// gc object task.
	DefaultGlobalBatchGcObjectTimeInterval int = 30 * 60
	// DefaultGlobalGcObjectBlockInterval defines the default blocks number for getting
	// deleted objects.
	DefaultGlobalGcObjectBlockInterval uint64 = 500
	// DefaultGlobalGcObjectSafeBlockDistance defines the default distance form current block
	// height to gc the deleted object.
	DefaultGlobalGcObjectSafeBlockDistance uint64 = 1000
	// DefaultGlobalSyncConsensusInfoInterval defines the default interval for sync the sp
	// info list to sp db.
	DefaultGlobalSyncConsensusInfoInterval uint64 = 2
	// DefaultStatisticsOutputInterval defines the default interval for output statistics info,
	// it is used to log and debug.
	DefaultStatisticsOutputInterval int = 60
)

func NewManageModular(app *gfspapp.GfSpBaseApp, cfg *gfspconfig.GfSpConfig) (coremodule.Modular, error) {
	manager := &ManageModular{baseApp: app}
	if err := DefaultManagerOptions(manager, cfg); err != nil {
		return nil, err
	}
	return manager, nil
}

func DefaultManagerOptions(manager *ManageModular, cfg *gfspconfig.GfSpConfig) error {
	if cfg.Parallel.GlobalMaxUploadingParallel == 0 {
		cfg.Parallel.GlobalMaxUploadingParallel = DefaultGlobalMaxUploadingNumber
	}
	if cfg.Parallel.GlobalUploadObjectParallel == 0 {
		cfg.Parallel.GlobalUploadObjectParallel = DefaultGlobalUploadObjectParallel
	}
	if cfg.Parallel.GlobalReplicatePieceParallel == 0 {
		cfg.Parallel.GlobalReplicatePieceParallel = DefaultGlobalReplicatePieceParallel
	}
	if cfg.Parallel.GlobalSealObjectParallel == 0 {
		cfg.Parallel.GlobalSealObjectParallel = DefaultGlobalSealObjectParallel
	}
	if cfg.Parallel.GlobalReceiveObjectParallel == 0 {
		cfg.Parallel.GlobalReceiveObjectParallel = DefaultGlobalReceiveObjectParallel
	}
	if cfg.Parallel.GlobalGCObjectParallel == 0 {
		cfg.Parallel.GlobalGCObjectParallel = DefaultGlobalGCObjectParallel
	}
	if cfg.Parallel.GlobalGCZombieParallel == 0 {
		cfg.Parallel.GlobalGCZombieParallel = DefaultGlobalGCZombieParallel
	}
	if cfg.Parallel.GlobalGCMetaParallel == 0 {
		cfg.Parallel.GlobalGCMetaParallel = DefaultGlobalGCMetaParallel
	}
	if cfg.Parallel.GlobalDownloadObjectTaskCacheSize == 0 {
		cfg.Parallel.GlobalDownloadObjectTaskCacheSize = DefaultGlobalDownloadObjectTaskCacheSize
	}
	if cfg.Parallel.GlobalChallengePieceTaskCacheSize == 0 {
		cfg.Parallel.GlobalChallengePieceTaskCacheSize = DefaultGlobalChallengePieceTaskCacheSize
	}
	if cfg.Parallel.GlobalBatchGcObjectTimeInterval == 0 {
		cfg.Parallel.GlobalBatchGcObjectTimeInterval = DefaultGlobalBatchGcObjectTimeInterval
	}
	if cfg.Parallel.GlobalGcObjectBlockInterval == 0 {
		cfg.Parallel.GlobalGcObjectBlockInterval = DefaultGlobalGcObjectBlockInterval
	}
	if cfg.Parallel.GlobalGcObjectSafeBlockDistance == 0 {
		cfg.Parallel.GlobalGcObjectSafeBlockDistance = DefaultGlobalGcObjectSafeBlockDistance
	}
	if cfg.Parallel.GlobalSyncConsensusInfoInterval == 0 {
		cfg.Parallel.GlobalSyncConsensusInfoInterval = DefaultGlobalSyncConsensusInfoInterval
	}
	manager.statisticsOutputInterval = DefaultStatisticsOutputInterval
	manager.maxUploadObjectNumber = cfg.Parallel.GlobalMaxUploadingParallel
	manager.gcObjectTimeInterval = cfg.Parallel.GlobalBatchGcObjectTimeInterval
	manager.gcObjectBlockInterval = cfg.Parallel.GlobalGcObjectBlockInterval
	manager.gcSafeBlockDistance = cfg.Parallel.GlobalGcObjectSafeBlockDistance
	manager.syncConsensusInfoInterval = cfg.Parallel.GlobalSyncConsensusInfoInterval
	manager.uploadQueue = cfg.Customize.NewStrategyTQueueFunc(
		manager.Name()+"-upload-object", cfg.Parallel.GlobalUploadObjectParallel)
	manager.replicateQueue = cfg.Customize.NewStrategyTQueueWithLimitFunc(
		manager.Name()+"-replicate-piece", cfg.Parallel.GlobalReplicatePieceParallel)
	manager.sealQueue = cfg.Customize.NewStrategyTQueueWithLimitFunc(
		manager.Name()+"-seal-object", cfg.Parallel.GlobalSealObjectParallel)
	manager.receiveQueue = cfg.Customize.NewStrategyTQueueWithLimitFunc(
		manager.Name()+"-confirm-receive-piece", cfg.Parallel.GlobalReceiveObjectParallel)
	manager.gcObjectQueue = cfg.Customize.NewStrategyTQueueWithLimitFunc(
		manager.Name()+"-gc-object", cfg.Parallel.GlobalGCObjectParallel)
	manager.gcZombieQueue = cfg.Customize.NewStrategyTQueueWithLimitFunc(
		manager.Name()+"-gc-zombie", cfg.Parallel.GlobalGCZombieParallel)
	manager.gcMetaQueue = cfg.Customize.NewStrategyTQueueWithLimitFunc(
		manager.Name()+"-gc-meta", cfg.Parallel.GlobalGCMetaParallel)
	manager.downloadQueue = cfg.Customize.NewStrategyTQueueFunc(
		manager.Name()+"-cache-download-object", cfg.Parallel.GlobalDownloadObjectTaskCacheSize)
	manager.challengeQueue = cfg.Customize.NewStrategyTQueueFunc(
		manager.Name()+"-cache-challenge-piece", cfg.Parallel.GlobalChallengePieceTaskCacheSize)
	return nil
}