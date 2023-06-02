package command

import (
	"context"
	"fmt"

	"github.com/bnb-chain/greenfield-storage-provider/base/types/gfsptask"
	"github.com/bnb-chain/greenfield-storage-provider/cmd/utils"
	"github.com/bnb-chain/greenfield-storage-provider/core/spdb"
	coretask "github.com/bnb-chain/greenfield-storage-provider/core/task"
	"github.com/bnb-chain/greenfield-storage-provider/util"
	storagetypes "github.com/bnb-chain/greenfield/x/storage/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/urfave/cli/v2"
)

const (
	DebugCommandPrefix = "gfsp-cli-debug-"
)

var DebugCreateBucketApprovalCmd = &cli.Command{
	Action: createBucketApproval,
	Name:   "debug.create.bucket.approval",
	Usage:  "Create random CreateBucketApproval and send to approver for debugging and testing",
	Flags: []cli.Flag{
		utils.ConfigFileFlag,
	},
	Category: "DEBUG COMMANDS",
	Description: `The debug.create.bucket.approval command create a random 
CreateBucketApproval request and send it to approver for debugging and testing
the approver on Dev Env.`,
}

// createBucketApproval is the debug.create.bucket.approval command action.
func createBucketApproval(ctx *cli.Context) error {
	cfg, err := utils.MakeConfig(ctx)
	if err != nil {
		return err
	}
	client := utils.MakeGfSpClient(cfg)

	msg := &storagetypes.MsgCreateBucket{
		BucketName:        DebugCommandPrefix + util.GetRandomBucketName(),
		PrimarySpApproval: &storagetypes.Approval{},
	}
	task := &gfsptask.GfSpCreateBucketApprovalTask{}
	task.InitApprovalCreateBucketTask(msg, coretask.UnSchedulingPriority)
	allow, res, err := client.AskCreateBucketApproval(context.Background(), task)
	if err != nil {
		return err
	}
	if !allow {
		return fmt.Errorf("refuse create bucket approval")
	}
	fmt.Printf("succeed to ask create bucket approval, BucketName: %s, ExpiredHeight: %d \n",
		res.GetCreateBucketInfo().GetBucketName(), res.GetCreateBucketInfo().GetPrimarySpApproval().GetExpiredHeight())
	return nil
}

var DebugCreateObjectApprovalCmd = &cli.Command{
	Action: createObjectApproval,
	Name:   "debug.create.object.approval",
	Usage:  "Create random CreateObjectApproval and send to approver for debugging and testing",
	Flags: []cli.Flag{
		utils.ConfigFileFlag,
	},
	Category: "DEBUG COMMANDS",
	Description: `The debug.create.object.approval command create a random 
CreateObjectApproval request and send it to approver for debugging and testing
the approver on Dev Env.`,
}

// createBucketApproval is the debug.create.bucket.approval command action.
func createObjectApproval(ctx *cli.Context) error {
	cfg, err := utils.MakeConfig(ctx)
	if err != nil {
		return err
	}
	client := utils.MakeGfSpClient(cfg)

	msg := &storagetypes.MsgCreateObject{
		BucketName:        DebugCommandPrefix + util.GetRandomBucketName(),
		ObjectName:        DebugCommandPrefix + util.GetRandomObjectName(),
		PrimarySpApproval: &storagetypes.Approval{},
	}
	task := &gfsptask.GfSpCreateObjectApprovalTask{}
	task.InitApprovalCreateObjectTask(msg, coretask.UnSchedulingPriority)
	allow, res, err := client.AskCreateObjectApproval(context.Background(), task)
	if err != nil {
		return err
	}
	if !allow {
		return fmt.Errorf("refuse create object approval")
	}
	fmt.Printf("succeed to ask create object approval, BucketName: %s, ObjectName: %s, ExpiredHeight: %d \n",
		res.GetCreateObjectInfo().GetBucketName(), res.GetCreateObjectInfo().GetBucketName(),
		res.GetCreateObjectInfo().GetPrimarySpApproval().GetExpiredHeight())
	return nil
}

var DebugReplicateApprovalCmd = &cli.Command{
	Action: replicatePieceApprovalAction,
	Name:   "debug.ask.replicate.approval",
	Usage:  "Create random ObjectInfo and send to p2p for debugging and testing p2p protocol network",
	Flags: []cli.Flag{
		utils.ConfigFileFlag,
		number,
	},
	Category: "DEBUG COMMANDS",
	Description: `The debug.ask.replicate.approval command create a random 
ObjectInfo and send it to p2p node for debugging and testing p2p protocol 
networkthe Dev Env.`,
}

// createBucketApproval is the debug.create.bucket.approval command action.
func replicatePieceApprovalAction(ctx *cli.Context) error {
	cfg, err := utils.MakeConfig(ctx)
	if err != nil {
		return err
	}
	client := utils.MakeGfSpClient(cfg)

	objectInfo := &storagetypes.ObjectInfo{
		Id:         sdk.NewUint(uint64(util.RandInt64(0, 100000))),
		BucketName: DebugCommandPrefix + util.GetRandomBucketName(),
		ObjectName: DebugCommandPrefix + util.GetRandomObjectName(),
	}
	task := &gfsptask.GfSpReplicatePieceApprovalTask{}
	task.InitApprovalReplicatePieceTask(objectInfo, &storagetypes.Params{},
		coretask.UnSchedulingPriority, GfSpCliUserName)

	expectNumber := ctx.Int(number.Name)
	approvals, err := client.AskSecondaryReplicatePieceApproval(
		context.Background(), task, expectNumber, expectNumber, 10)
	if err != nil {
		return err
	}
	fmt.Printf("receive %d accepted approvals\n", len(approvals))

	db, err := utils.MakeSPDB(cfg)
	if err != nil {
		return err
	}
	for _, approval := range approvals {
		spInfo, err := db.GetSpByAddress(approval.GetApprovedSpOperatorAddress(), spdb.OperatorAddressType)
		if err != nil {
			return err
		}
		fmt.Printf("%s[%s] accepted\n", approval.GetApprovedSpOperatorAddress(), spInfo.GetEndpoint())
	}
	return nil
}
