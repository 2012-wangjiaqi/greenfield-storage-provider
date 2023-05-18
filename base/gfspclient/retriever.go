package gfspclient

import (
	"context"

	"google.golang.org/grpc"

	retrievertypes "github.com/bnb-chain/greenfield-storage-provider/modular/retriever/types"
	metatypes "github.com/bnb-chain/greenfield-storage-provider/service/metadata/types"
	storagetypes "github.com/bnb-chain/greenfield/x/storage/types"
)

func (s *GfSpClient) GetAccountBucketNumber(
	ctx context.Context,
	account string,
	opts ...grpc.DialOption) (int64, error) {
	conn, err := s.Connection(ctx, s.metadataEndpoint, opts...)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	req := &metatypes.GetUserBucketsCountRequest{
		AccountId: account,
	}
	resp, err := metatypes.NewMetadataServiceClient(conn).GetUserBucketsCount(ctx, req)
	if err != nil {
		return 0, ErrRpcUnknown
	}
	return resp.GetCount(), nil
}

func (s *GfSpClient) ListDeletedObjectsByBlockNumberRange(
	ctx context.Context,
	spOperatorAddress string,
	startBlockNumber uint64,
	endBlockNumber uint64,
	includePrivate bool,
	opts ...grpc.DialOption) ([]*metatypes.Object, uint64, error) {
	conn, err := s.Connection(ctx, s.metadataEndpoint, opts...)
	if err != nil {
		return nil, uint64(0), err
	}
	defer conn.Close()
	req := &metatypes.ListDeletedObjectsByBlockNumberRangeRequest{
		StartBlockNumber: int64(startBlockNumber),
		EndBlockNumber:   int64(endBlockNumber),
		IsFullList:       includePrivate,
	}
	resp, err := metatypes.NewMetadataServiceClient(conn).ListDeletedObjectsByBlockNumberRange(ctx, req)
	if err != nil {
		return nil, uint64(0), ErrRpcUnknown
	}
	return resp.GetObjects(), uint64(resp.GetEndBlockNumber()), nil
}

func (s *GfSpClient) GetUserBuckets(
	ctx context.Context,
	account string,
	opts ...grpc.DialOption) ([]*retrievertypes.Bucket, error) {
	conn, err := s.Connection(ctx, s.metadataEndpoint, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	req := &retrievertypes.GfSpGetUserBucketsRequest{
		AccountId: account,
	}
	resp, err := retrievertypes.NewGfSpRetrieverServiceClient(conn).GfSpGetUserBuckets(ctx, req)
	if err != nil {
		return nil, ErrRpcUnknown
	}
	return resp.GetBuckets(), nil
}

func (s *GfSpClient) GetBucketReadQuota(
	ctx context.Context,
	bucket *storagetypes.BucketInfo,
	yearMonth string,
	opts ...grpc.DialOption) (
	uint64, uint64, uint64, error) {
	conn, err := s.Connection(ctx, s.retrieverEndpoint, opts...)
	if err != nil {
		return uint64(0), uint64(0), uint64(0), err
	}
	defer conn.Close()
	req := &retrievertypes.GfSpGetBucketReadQuotaRequest{
		BucketInfo: bucket,
		YearMonth:  yearMonth,
	}
	resp, err := retrievertypes.NewGfSpRetrieverServiceClient(conn).GfSpGetBucketReadQuota(ctx, req)
	if err != nil {
		return uint64(0), uint64(0), uint64(0), ErrRpcUnknown
	}
	return resp.GetChargedQuotaSize(), resp.GetSpFreeQuotaSize(), resp.GetConsumedSize(), resp.GetErr()
}

func (s *GfSpClient) ListBucketReadRecord(
	ctx context.Context,
	bucket *storagetypes.BucketInfo,
	startTimestampUs, endTimestampUs, maxRecordNum int64,
	opts ...grpc.DialOption) (
	[]*retrievertypes.ReadRecord, int64, error) {
	conn, err := s.Connection(ctx, s.retrieverEndpoint, opts...)
	if err != nil {
		return nil, 0, err
	}
	defer conn.Close()
	req := &retrievertypes.GfSpListBucketReadRecordRequest{
		BucketInfo:       bucket,
		StartTimestampUs: startTimestampUs,
		EndTimestampUs:   endTimestampUs,
		MaxRecordNum:     maxRecordNum,
	}
	resp, err := retrievertypes.NewGfSpRetrieverServiceClient(conn).GfSpListBucketReadRecord(ctx, req)
	if err != nil {
		return nil, 0, ErrRpcUnknown
	}
	return resp.GetReadRecords(), resp.GetNextStartTimestampUs(), resp.GetErr()
}

func (s *GfSpClient) GetUploadObjectState(
	ctx context.Context,
	objectID uint64,
	opts ...grpc.DialOption) (int32, error) {
	conn, err := s.Connection(ctx, s.retrieverEndpoint, opts...)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	req := &retrievertypes.GfSpQueryUploadProgressRequest{
		ObjectId: objectID,
	}
	resp, err := retrievertypes.NewGfSpRetrieverServiceClient(conn).GfSpQueryUploadProgress(ctx, req)
	if err != nil {
		return 0, ErrRpcUnknown
	}
	return int32(resp.GetState()), err
}