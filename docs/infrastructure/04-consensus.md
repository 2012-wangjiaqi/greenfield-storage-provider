# Consensus 

Consensus is the basic component that access the consensus data that products from greenfield.
However, considering that these consensus data may undergo secondary processing in optimization
scenarios, SP abstracts the interface for consensus data access, so that SP deployers can 
obtain consensus data from different data sources, such as: greenfield validator, greenfield 
full node or other greenfield off-chain data services.


## Consensus Interface
```go
// Consensus is the interface to query greenfield consensus data. the consensus
// data can come from validator, full-node, or other off-chain data service
type Consensus interface {
	// CurrentHeight returns the current greenfield height - 1,
	CurrentHeight(ctx context.Context) (uint64, error)
	// HasAccount returns an indicator whether the account has been created.
	HasAccount(ctx context.Context, account string) (bool, error)
	// QuerySPInfo returns all SP info.
	QuerySPInfo(ctx context.Context) ([]*sptypes.StorageProvider, error)
	// QueryStorageParams returns the storage params.
	QueryStorageParams(ctx context.Context) (params *storagetypes.Params, err error)
	// QueryBucketInfo returns the bucket info by bucket name.
	QueryBucketInfo(ctx context.Context, bucket string) (*storagetypes.BucketInfo, error)
	// QueryObjectInfo returns the object info by bucket and object name.
	QueryObjectInfo(ctx context.Context, bucket, object string) (*storagetypes.ObjectInfo, error)
	// QueryObjectInfoByID returns the object info by object ID.
	QueryObjectInfoByID(ctx context.Context, objectID string) (*storagetypes.ObjectInfo, error)
	// QueryBucketInfoAndObjectInfo returns the bucket and object info by bucket and object name.
	QueryBucketInfoAndObjectInfo(ctx context.Context, bucket, object string) (*storagetypes.BucketInfo, *storagetypes.ObjectInfo, error)
	// QueryPaymentStreamRecord returns the account payment status.
	QueryPaymentStreamRecord(ctx context.Context, account string) (*paymenttypes.StreamRecord, error)
	// VerifyGetObjectPermission returns an indicator whether the account has permission to get object.
	VerifyGetObjectPermission(ctx context.Context, account, bucket, object string) (bool, error)
	// VerifyPutObjectPermission returns an indicator whether the account has permission to put object.
	VerifyPutObjectPermission(ctx context.Context, account, bucket, object string) (bool, error)
	// ListenObjectSeal returns an indicator whether the object is successfully sealed before timeOutHeight.
	ListenObjectSeal(ctx context.Context, objectID uint64, timeOutHeight int) (bool, error)
	// Close the Consensus interface.
	Close() error
}

```