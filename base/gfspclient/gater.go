package gfspclient

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bnb-chain/greenfield-storage-provider/base/types/gfsptask"
	coretask "github.com/bnb-chain/greenfield-storage-provider/core/task"
	"github.com/bnb-chain/greenfield-storage-provider/model"
	"github.com/bnb-chain/greenfield-storage-provider/pkg/log"
)

func (s *GfSpClient) ReplicatePieceToSecondary(
	ctx context.Context,
	endpoint string,
	approval coretask.ApprovalReplicatePieceTask,
	receive coretask.ReceivePieceTask,
	data []byte) error {
	req, err := http.NewRequest(http.MethodPut, endpoint+model.ReplicateObjectPiecePath, bytes.NewReader(data))
	if err != nil {
		log.CtxErrorw(ctx, "client failed to connect gateway", "endpoint", endpoint, "error", err)
		return err
	}
	approvalTask := approval.(*gfsptask.GfSpReplicatePieceApprovalTask)
	approvalMsg, err := json.Marshal(approvalTask)
	if err != nil {
		return err
	}
	approvalHeader := hex.EncodeToString(approvalMsg)

	receiveTask := receive.(*gfsptask.GfSpReceivePieceTask)
	receiveMsg, err := json.Marshal(receiveTask)
	if err != nil {
		return err
	}
	receiveHeader := hex.EncodeToString(receiveMsg)
	req.Header.Add(model.GnfdReplicatePieceApprovalHeader, approvalHeader)
	req.Header.Add(model.GnfdReceiveMsgHeader, receiveHeader)
	resp, err := s.HttpClient(ctx).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to replicate piece, StatusCode(%d) Endpoint(%s)", resp.StatusCode, endpoint)
	}
	return nil
}

func (s *GfSpClient) DoneReplicatePieceToSecondary(
	ctx context.Context,
	endpoint string,
	approval coretask.ApprovalReplicatePieceTask,
	receive coretask.ReceivePieceTask,
) ([]byte, []byte, error) {
	req, err := http.NewRequest(http.MethodPut, endpoint+model.ReplicateObjectPiecePath, nil)
	if err != nil {
		log.CtxErrorw(ctx, "client failed to connect gateway", "endpoint", endpoint, "error", err)
		return nil, nil, err
	}
	approvalTask := approval.(*gfsptask.GfSpReplicatePieceApprovalTask)
	approvalMsg, err := json.Marshal(approvalTask)
	if err != nil {
		return nil, nil, err
	}
	approvalHeader := hex.EncodeToString(approvalMsg)

	receiveTask := receive.(*gfsptask.GfSpReceivePieceTask)
	receiveMsg, err := json.Marshal(receiveTask)
	if err != nil {
		return nil, nil, err
	}
	receiveHeader := hex.EncodeToString(receiveMsg)
	req.Header.Add(model.GnfdReplicatePieceApprovalHeader, approvalHeader)
	req.Header.Add(model.GnfdReceiveMsgHeader, receiveHeader)
	resp, err := s.HttpClient(ctx).Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("failed to replicate piece, StatusCode(%d)", resp.StatusCode)
	}
	integrity, err := hex.DecodeString(resp.Header.Get(model.GnfdIntegrityHashHeader))
	if err != nil {
		return nil, nil, err
	}
	signature, err := hex.DecodeString(resp.Header.Get(model.GnfdIntegrityHashSignatureHeader))
	if err != nil {
		return nil, nil, err
	}
	return integrity, signature, nil
}
