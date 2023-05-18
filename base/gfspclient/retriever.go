package gfspclient

import (
	"context"
	"github.com/bnb-chain/greenfield-storage-provider/pkg/log"
	"google.golang.org/grpc"

	"github.com/bnb-chain/greenfield-storage-provider/modular/retriever/types"
	payment_types "github.com/bnb-chain/greenfield/x/payment/types"
	permission_types "github.com/bnb-chain/greenfield/x/permission/types"
	storage_types "github.com/bnb-chain/greenfield/x/storage/types"
)

func (s *GfSpClient) GetUserBucketsCount(
	ctx context.Context,
	account string,
	opts ...grpc.DialOption) (int64, error) {
	conn, err := s.Connection(ctx, s.metadataEndpoint, opts...)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	req := &types.GfSpGetUserBucketsCountRequest{
		AccountId: account,
	}
	resp, err := types.NewGfSpRetrieverServiceClient(conn).GfSpGetUserBucketsCount(ctx, req)
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
	opts ...grpc.DialOption) ([]*types.Object, uint64, error) {
	conn, err := s.Connection(ctx, s.metadataEndpoint, opts...)
	if err != nil {
		return nil, uint64(0), err
	}
	defer conn.Close()
	req := &types.GfSpListDeletedObjectsByBlockNumberRangeRequest{
		StartBlockNumber: int64(startBlockNumber),
		EndBlockNumber:   int64(endBlockNumber),
		IsFullList:       includePrivate,
	}
	resp, err := types.NewGfSpRetrieverServiceClient(conn).GfSpListDeletedObjectsByBlockNumberRange(ctx, req)
	if err != nil {
		return nil, uint64(0), ErrRpcUnknown
	}
	return resp.GetObjects(), uint64(resp.GetEndBlockNumber()), nil
}

func (s *GfSpClient) GetUserBuckets(
	ctx context.Context,
	account string,
	opts ...grpc.DialOption) ([]*types.Bucket, error) {
	conn, err := s.Connection(ctx, s.metadataEndpoint, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	req := &types.GfSpGetUserBucketsRequest{
		AccountId: account,
	}
	resp, err := types.NewGfSpRetrieverServiceClient(conn).GfSpGetUserBuckets(ctx, req)
	if err != nil {
		return nil, ErrRpcUnknown
	}
	return resp.GetBuckets(), nil
}

// ListObjectsByBucketName list objects info by a bucket name
func (s *GfSpClient) ListObjectsByBucketName(
	ctx context.Context,
	bucketName string,
	accountId string,
	maxKeys uint64,
	startAfter string,
	continuationToken string,
	delimiter string,
	prefix string,
	opts ...grpc.DialOption) (
	objects []*types.Object,
	KeyCount uint64,
	MaxKeys uint64,
	IsTruncated bool,
	NextContinuationToken string,
	Name string,
	Prefix string,
	Delimiter string,
	CommonPrefixes []string,
	ContinuationToken string,
	err error) {
	conn, err := s.Connection(ctx, s.metadataEndpoint, opts...)
	if err != nil {
		return nil, 0, 0, false, "", "", "", "", nil, "", err
	}
	defer conn.Close()

	req := &types.GfSpListObjectsByBucketNameRequest{
		BucketName:        bucketName,
		AccountId:         accountId,
		MaxKeys:           maxKeys,
		StartAfter:        startAfter,
		ContinuationToken: continuationToken,
		Delimiter:         delimiter,
		Prefix:            prefix,
	}

	resp, err := types.NewGfSpRetrieverServiceClient(conn).GfSpListObjectsByBucketName(ctx, req)
	ctx = log.Context(ctx, resp)
	if err != nil {
		log.CtxErrorw(ctx, "failed to send list objects by bucket name rpc", "error", err)
		return nil, 0, 0, false, "", "", "", "", nil, "", err
	}
	return resp.GetObjects(), resp.GetKeyCount(), resp.GetMaxKeys(), resp.GetIsTruncated(), resp.GetNextContinuationToken(),
		resp.GetName(), resp.GetPrefix(), resp.GetDelimiter(), resp.GetCommonPrefixes(), resp.GetContinuationToken(), nil
}

// GetBucketByBucketName get bucket info by a bucket name
func (s *GfSpClient) GetBucketByBucketName(
	ctx context.Context,
	bucketName string,
	isFullList bool,
	opts ...grpc.DialOption) (*types.Bucket, error) {
	conn, err := s.Connection(ctx, s.metadataEndpoint, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req := &types.GfSpGetBucketByBucketNameRequest{
		BucketName: bucketName,
		IsFullList: isFullList,
	}

	resp, err := types.NewGfSpRetrieverServiceClient(conn).GfSpGetBucketByBucketName(ctx, req)
	ctx = log.Context(ctx, resp)
	if err != nil {
		log.CtxErrorw(ctx, "failed to send get bucket rpc by bucket name", "error", err)
		return nil, err
	}
	return resp.GetBucket(), nil
}

// GetBucketByBucketID get bucket info by a bucket id
func (s *GfSpClient) GetBucketByBucketID(ctx context.Context,
	bucketId int64,
	isFullList bool,
	opts ...grpc.DialOption) (*types.GfSpGetBucketByBucketIDResponse, error) {
	conn, err := s.Connection(ctx, s.metadataEndpoint, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req := &types.GfSpGetBucketByBucketIDRequest{
		BucketId:   bucketId,
		IsFullList: isFullList,
	}

	resp, err := types.NewGfSpRetrieverServiceClient(conn).GfSpGetBucketByBucketID(ctx, req)
	ctx = log.Context(ctx, resp)
	if err != nil {
		log.CtxErrorw(ctx, "failed to send get bucket by bucket id rpc", "error", err)
		return nil, err
	}
	return resp, nil
}

// ListExpiredBucketsBySp list buckets that are expired by specific sp
func (s *GfSpClient) ListExpiredBucketsBySp(ctx context.Context, createAt int64, primarySpAddress string, limit int64, opts ...grpc.DialOption) ([]*types.Bucket, error) {
	conn, err := s.Connection(ctx, s.metadataEndpoint, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req := &types.GfSpListExpiredBucketsBySpRequest{
		CreateAt:         createAt,
		PrimarySpAddress: primarySpAddress,
		Limit:            limit,
	}

	resp, err := types.NewGfSpRetrieverServiceClient(conn).GfSpListExpiredBucketsBySp(ctx, req)
	ctx = log.Context(ctx, resp)
	if err != nil {
		log.CtxErrorw(ctx, "failed to send list expired buckets by sp rpc", "error", err)
		return nil, err
	}
	return resp.GetBuckets(), nil
}

// GetObjectMeta get object metadata
func (s *GfSpClient) GetObjectMeta(
	ctx context.Context,
	objectName string,
	bucketName string,
	isFullList bool,
	opts ...grpc.DialOption) (*types.Object, error) {
	conn, err := s.Connection(ctx, s.metadataEndpoint, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req := &types.GfSpGetObjectMetaRequest{
		ObjectName: objectName,
		BucketName: bucketName,
		IsFullList: isFullList,
	}

	resp, err := types.NewGfSpRetrieverServiceClient(conn).GfSpGetObjectMeta(ctx, req)
	ctx = log.Context(ctx, resp)
	if err != nil {
		log.CtxErrorw(ctx, "failed to send get object meta rpc", "error", err)
		return nil, err
	}
	return resp.GetObject(), nil
}

// GetPaymentByBucketName get bucket payment info by a bucket name
func (s *GfSpClient) GetPaymentByBucketName(
	ctx context.Context,
	bucketName string,
	isFullList bool,
	opts ...grpc.DialOption) (*payment_types.StreamRecord, error) {
	conn, err := s.Connection(ctx, s.metadataEndpoint, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req := &types.GfSpGetPaymentByBucketNameRequest{
		BucketName: bucketName,
		IsFullList: isFullList,
	}

	resp, err := types.NewGfSpRetrieverServiceClient(conn).GfSpGetPaymentByBucketName(ctx, req)
	ctx = log.Context(ctx, resp)
	if err != nil {
		log.CtxErrorw(ctx, "failed to send get payment by bucket name rpc", "error", err)
		return nil, err
	}
	return resp.GetStreamRecord(), nil
}

// GetPaymentByBucketID get bucket payment info by a bucket id
func (s *GfSpClient) GetPaymentByBucketID(
	ctx context.Context,
	bucketID int64,
	isFullList bool,
	opts ...grpc.DialOption) (*payment_types.StreamRecord, error) {
	conn, err := s.Connection(ctx, s.metadataEndpoint, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req := &types.GfSpGetPaymentByBucketIDRequest{
		BucketId:   bucketID,
		IsFullList: isFullList,
	}

	resp, err := types.NewGfSpRetrieverServiceClient(conn).GfSpGetPaymentByBucketID(ctx, req)
	ctx = log.Context(ctx, resp)
	if err != nil {
		log.CtxErrorw(ctx, "failed to send get payment by bucket id rpc", "error", err)
		return nil, err
	}
	return resp.GetStreamRecord(), nil
}

// VerifyPermission Verify the input account’s permission to input items
func (s *GfSpClient) VerifyPermission(
	ctx context.Context,
	Operator string,
	bucketName string,
	objectName string,
	actionType permission_types.ActionType,
	opts ...grpc.DialOption) (*permission_types.Effect, error) {
	conn, err := s.Connection(ctx, s.metadataEndpoint, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req := &storage_types.QueryVerifyPermissionRequest{
		Operator:   Operator,
		BucketName: bucketName,
		ObjectName: objectName,
		ActionType: actionType,
	}

	resp, err := types.NewGfSpRetrieverServiceClient(conn).GfSpVerifyPermission(ctx, req)
	ctx = log.Context(ctx, resp)
	if err != nil {
		log.CtxErrorw(ctx, "failed to send verify permission rpc", "error", err)
		return nil, err
	}
	return &resp.Effect, nil
}

// GetBucketMeta get bucket info along with its related info such as payment
func (s *GfSpClient) GetBucketMeta(
	ctx context.Context,
	bucketName string,
	isFullList bool,
	opts ...grpc.DialOption) (*types.Bucket, *payment_types.StreamRecord, error) {
	conn, err := s.Connection(ctx, s.metadataEndpoint, opts...)
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	req := &types.GfSpGetBucketMetaRequest{
		BucketName: bucketName,
		IsFullList: isFullList,
	}

	resp, err := types.NewGfSpRetrieverServiceClient(conn).GfSpGetBucketMeta(ctx, req)
	ctx = log.Context(ctx, resp)
	if err != nil {
		log.CtxErrorw(ctx, "failed to send get bucket meta rpc", "error", err)
		return nil, nil, err
	}
	return resp.GetBucket(), resp.GetStreamRecord(), nil
}

func (s *GfSpClient) GetBucketReadQuota(
	ctx context.Context,
	bucket *storage_types.BucketInfo,
	yearMonth string,
	opts ...grpc.DialOption) (
	uint64, uint64, uint64, error) {
	conn, err := s.Connection(ctx, s.retrieverEndpoint, opts...)
	if err != nil {
		return uint64(0), uint64(0), uint64(0), err
	}
	defer conn.Close()
	req := &types.GfSpGetBucketReadQuotaRequest{
		BucketInfo: bucket,
		YearMonth:  yearMonth,
	}
	resp, err := types.NewGfSpRetrieverServiceClient(conn).GfSpGetBucketReadQuota(ctx, req)
	if err != nil {
		return uint64(0), uint64(0), uint64(0), ErrRpcUnknown
	}
	return resp.GetChargedQuotaSize(), resp.GetSpFreeQuotaSize(), resp.GetConsumedSize(), resp.GetErr()
}

func (s *GfSpClient) ListBucketReadRecord(
	ctx context.Context,
	bucket *storage_types.BucketInfo,
	startTimestampUs, endTimestampUs, maxRecordNum int64,
	opts ...grpc.DialOption) (
	[]*types.ReadRecord, int64, error) {
	conn, err := s.Connection(ctx, s.retrieverEndpoint, opts...)
	if err != nil {
		return nil, 0, err
	}
	defer conn.Close()
	req := &types.GfSpListBucketReadRecordRequest{
		BucketInfo:       bucket,
		StartTimestampUs: startTimestampUs,
		EndTimestampUs:   endTimestampUs,
		MaxRecordNum:     maxRecordNum,
	}
	resp, err := types.NewGfSpRetrieverServiceClient(conn).GfSpListBucketReadRecord(ctx, req)
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
	req := &types.GfSpQueryUploadProgressRequest{
		ObjectId: objectID,
	}
	resp, err := types.NewGfSpRetrieverServiceClient(conn).GfSpQueryUploadProgress(ctx, req)
	if err != nil {
		return 0, ErrRpcUnknown
	}
	return int32(resp.GetState()), err
}
