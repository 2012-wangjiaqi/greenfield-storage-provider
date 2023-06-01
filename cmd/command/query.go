package command

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/bnb-chain/greenfield-common/go/hash"
	"github.com/bnb-chain/greenfield-storage-provider/base/types/gfsptask"
	coretask "github.com/bnb-chain/greenfield-storage-provider/core/task"
	"github.com/urfave/cli/v2"

	"github.com/bnb-chain/greenfield-storage-provider/base/gfspapp"
	"github.com/bnb-chain/greenfield-storage-provider/base/gfspclient"
	"github.com/bnb-chain/greenfield-storage-provider/base/gfspconfig"
	"github.com/bnb-chain/greenfield-storage-provider/cmd/utils"
	"github.com/bnb-chain/greenfield-storage-provider/pkg/log"
)

const (
	GfSpCliUserName = "gfsp-cli"
)

var endpointFlag = &cli.StringFlag{
	Name:  "n",
	Usage: "The address of machine that to query tasks",
	Value: "",
}

var keyFlag = &cli.StringFlag{
	Name:  "k",
	Usage: "The sub key of task",
	Value: "",
}

var QueryTaskCmd = &cli.Command{
	Action:   queryTasksAction,
	Name:     "query.task",
	Usage:    "Query running tasks in modules by task sub key",
	Category: "QUERY COMMANDS",
	Flags: []cli.Flag{
		utils.ConfigFileFlag,
		endpointFlag,
		keyFlag,
	},
	Description: `Query running tasks in modules by task sub key, 
show the tasks that task key contains the inout key detail info`,
}

func queryTasksAction(ctx *cli.Context) error {
	endpoint := gfspapp.DefaultGrpcAddress
	if ctx.IsSet(utils.ConfigFileFlag.Name) {
		cfg := &gfspconfig.GfSpConfig{}
		err := utils.LoadConfig(ctx.String(utils.ConfigFileFlag.Name), cfg)
		if err != nil {
			log.Errorw("failed to load config file", "error", err)
			return err
		}
		endpoint = cfg.GrpcAddress
	}
	if ctx.IsSet(endpointFlag.Name) {
		endpoint = ctx.String(endpointFlag.Name)
	}
	if !ctx.IsSet(keyFlag.Name) {
		return fmt.Errorf("query key should be set")
	}
	key := ctx.String(keyFlag.Name)
	if len(key) == 0 {
		return fmt.Errorf("query key can not empty")
	}
	client := &gfspclient.GfSpClient{}
	infos, err := client.QueryTasks(context.Background(), endpoint, key)
	if err != nil {
		return err
	}
	if len(infos) == 0 {
		return fmt.Errorf("no task match the query key")
	}
	for _, info := range infos {
		fmt.Printf(info + "\n")
	}
	return nil
}

var objectIDFlag = &cli.StringFlag{
	Name:     "i",
	Usage:    "The ID key of Object",
	Required: true,
}

var GetObjectCmd = &cli.Command{
	Action: getObjectAction,
	Name:   "get.object",
	Usage:  "Get object payload data",
	Flags: []cli.Flag{
		utils.ConfigFileFlag,
		objectIDFlag,
	},
	Category: "QUERY COMMANDS",
	Description: `The get.object command send rpc request to downloader 
server to get object payload data`,
}

func getObjectAction(ctx *cli.Context) error {
	cfg, err := utils.MakeConfig(ctx)
	if err != nil {
		return err
	}
	client := utils.MakeGfSpClient(cfg)
	chain, err := utils.MakeGnfd(cfg)
	if err != nil {
		return err
	}

	objectID := ctx.String(objectIDFlag.Name)
	objectInfo, err := chain.QueryObjectInfoByID(context.Background(), objectID)
	if err != nil {
		return fmt.Errorf("failed to query object info, error: %v", err)
	}
	bucketInfo, err := chain.QueryBucketInfo(context.Background(), objectInfo.GetBucketName())
	if err != nil {
		return fmt.Errorf("failed to query object info, error: %v", err)
	}
	params, err := chain.QueryStorageParamsByTimestamp(context.Background(), objectInfo.GetCreateAt())
	if err != nil {
		return fmt.Errorf("failed to query storage params, error: %v", err)
	}
	task := &gfsptask.GfSpDownloadObjectTask{}
	task.InitDownloadObjectTask(objectInfo, bucketInfo, params, coretask.UnSchedulingPriority,
		GfSpCliUserName, 0, int64(objectInfo.GetPayloadSize()-1), 0, 0)
	data, err := client.GetObject(context.Background(), task)
	if err != nil {
		return fmt.Errorf("failed to get object, error: %v", err)
	}
	if err = os.WriteFile("./"+objectInfo.GetObjectName(), data, os.ModePerm); err != nil {
		fmt.Printf("failed to create file to wirte object payload data, error: %v", err)
	}
	fmt.Printf("succeed to get object\n BucketInfo: %s\n ObjectInfo: %s \nStorageParam: %s\n",
		bucketInfo.String(), objectInfo.String(), params.String())
	return nil
}

var replicateIdxFlag = &cli.Int64Flag{
	Name:     "r",
	Usage:    "The object replicate index of SP",
	Required: true,
}

var segmentIdxFlag = &cli.Uint64Flag{
	Name:     "s",
	Usage:    "The segment index",
	Required: true,
}

var ChallengePieceCmd = &cli.Command{
	Action: challengePieceAction,
	Name:   "challenge.piece",
	Usage:  "Challenge piece integrity hash",
	Flags: []cli.Flag{
		utils.ConfigFileFlag,
		objectIDFlag,
		replicateIdxFlag,
		segmentIdxFlag,
	},
	Category: "QUERY COMMANDS",
	Description: `The challenge.piece command send rpc request to downloader 
get integrity meta and check the piece checksums`,
}

func challengePieceAction(ctx *cli.Context) error {
	cfg, err := utils.MakeConfig(ctx)
	if err != nil {
		return err
	}
	client := utils.MakeGfSpClient(cfg)
	chain, err := utils.MakeGnfd(cfg)
	if err != nil {
		return err
	}

	objectID := ctx.String(objectIDFlag.Name)
	objectInfo, err := chain.QueryObjectInfoByID(context.Background(), objectID)
	if err != nil {
		return fmt.Errorf("failed to query object info, error: %v", err)
	}
	bucketInfo, err := chain.QueryBucketInfo(context.Background(), objectInfo.GetBucketName())
	if err != nil {
		return fmt.Errorf("failed to query object info, error: %v", err)
	}
	params, err := chain.QueryStorageParamsByTimestamp(context.Background(), objectInfo.GetCreateAt())
	if err != nil {
		return fmt.Errorf("failed to query storage params, error: %v", err)
	}

	replicateIdx := ctx.Int64(replicateIdxFlag.Name)
	segmentIdx := ctx.Uint64(segmentIdxFlag.Name)

	task := &gfsptask.GfSpChallengePieceTask{}
	task.InitChallengePieceTask(objectInfo, bucketInfo, params, coretask.UnSchedulingPriority,
		GfSpCliUserName, int32(replicateIdx), uint32(segmentIdx), 0, 0)
	integrityHash, checksums, data, err := client.GetChallengeInfo(context.Background(), task)
	if err != nil {
		return fmt.Errorf("failed to get challeneg info, error: %v", err)
	}
	fmt.Printf("integrity meta info: \n\n")
	fmt.Printf("integrity hash[%s]\n\n", hex.EncodeToString(integrityHash))
	for i, checksum := range checksums {
		fmt.Printf("piece[%d] checksum[%s]\n", i, hex.EncodeToString(checksum))
	}
	challengePieceChecksum := hash.GenerateChecksum(data)
	fmt.Printf("\nchallenge piece info: replicate_idx[%d], segment_idx[%d], piece_checksum[%s]\n\n",
		replicateIdx, segmentIdx, hex.EncodeToString(challengePieceChecksum))

	if !bytes.Equal(challengePieceChecksum, checksums[segmentIdx]) {
		return fmt.Errorf("piece data hash[%s] not equal to checksum list value[%s]",
			hex.EncodeToString(challengePieceChecksum), hex.EncodeToString(checksums[segmentIdx]))
	}

	if !bytes.Equal(integrityHash, hash.GenerateIntegrityHash(checksums)) {
		return fmt.Errorf("integrity hash[%s] mismatch checksum list result[%s]",
			hex.EncodeToString(integrityHash),
			hex.EncodeToString(hash.GenerateIntegrityHash(checksums)))
	}
	fmt.Printf("succeed to check integrity hash")
	return nil
}
