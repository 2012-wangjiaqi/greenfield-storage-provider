package command

import (
	"context"
	"fmt"

	"github.com/bnb-chain/greenfield-storage-provider/base/types/gfsptask"
	"github.com/bnb-chain/greenfield-storage-provider/cmd/utils"
	coretask "github.com/bnb-chain/greenfield-storage-provider/core/task"
	"github.com/bnb-chain/greenfield-storage-provider/util"
	storagetypes "github.com/bnb-chain/greenfield/x/storage/types"
	"github.com/urfave/cli/v2"
)

var DebugCreateBucketApprovalCmd = &cli.Command{
	Action:   createBucketApproval,
	Name:     "debug.create.bucket.approval",
	Usage:    "Create random CreateBucketApproval and send to approver for debugging and testing",
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
		BucketName: util.GetRandomBucketName(),
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
	fmt.Printf("succeed to ask create bucket approval, BucketName: %s, ExpiredHeight: %d",
		res.GetCreateBucketInfo().GetBucketName(), res.GetCreateBucketInfo().GetPrimarySpApproval().GetExpiredHeight())
	return nil
}

var DebugCreateObjectApprovalCmd = &cli.Command{
	Action:   createObjectApproval,
	Name:     "debug.create.object.approval",
	Usage:    "Create random CreateObjectApproval and send to approver for debugging and testing",
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
		ObjectName: util.GetRandomObjectName(),
		BucketName: util.GetRandomBucketName(),
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
	fmt.Printf("succeed to ask create object approval, BucketName: %s, ObjectName: %s, ExpiredHeight: %d",
		res.GetCreateObjectInfo().GetBucketName(), res.GetCreateObjectInfo().GetBucketName(),
		res.GetCreateObjectInfo().GetPrimarySpApproval().GetExpiredHeight())
	return nil
}
