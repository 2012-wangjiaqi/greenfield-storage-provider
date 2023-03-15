package service

import (
	"context"

	"cosmossdk.io/math"

	model "github.com/bnb-chain/greenfield-storage-provider/model/metadata"
	"github.com/bnb-chain/greenfield-storage-provider/pkg/log"
	metatypes "github.com/bnb-chain/greenfield-storage-provider/service/metadata/types"
	"github.com/bnb-chain/greenfield/x/storage/types"
)

// GetUserBuckets get buckets info by a user address
func (metadata *Metadata) GetUserBuckets(ctx context.Context, req *metatypes.MetadataServiceGetUserBucketsRequest) (resp *metatypes.MetadataServiceGetUserBucketsResponse, err error) {
	ctx = log.Context(ctx, req)
	defer func() {
		if err != nil {
			log.CtxErrorw(ctx, "failed to get user buckets", "err", err)
		} else {
			log.CtxInfow(ctx, "succeed to get user buckets")
		}
	}()

	var buckets []*model.Bucket
	//TODO:: cancel mock after impl db
	bucket1 := &model.Bucket{
		Owner:            "46765cbc-d30c-4f4a-a814-b68181fcab12",
		BucketName:       "BBC News",
		IsPublic:         true,
		ID:               "1",
		SourceType:       1,
		CreateAt:         1676530547,
		PaymentAddress:   "0x000000006b4BD0274e8f943201A922295D13fc28",
		PrimarySpAddress: "0x000000006b4BD0274e8f943201A922295D13fc28",
		ReadQuota:        2,
		PaymentPriceTime: 1677143663461,
	}
	bucket2 := &model.Bucket{
		Owner:            "46765cbc-d30c-4f4a-a814-b68181fcab12",
		BucketName:       "bnb-chain",
		IsPublic:         true,
		ID:               "2",
		SourceType:       1,
		CreateAt:         1676530547,
		PaymentAddress:   "0xF9A8db17431DD8563747D6FC770297E438Aa12eB",
		PrimarySpAddress: "0xF9A8db17431DD8563747D6FC770297E438Aa12eB",
		ReadQuota:        5,
		PaymentPriceTime: 1677143663461,
	}
	buckets = append(buckets, bucket1, bucket2)
	res := make([]*metatypes.Bucket, 0)

	for _, bucket := range buckets {
		res = append(res, &metatypes.Bucket{
			BucketInfo: &types.BucketInfo{
				Owner:            bucket.Owner,
				BucketName:       bucket.BucketName,
				IsPublic:         false,
				Id:               math.NewUintFromString(bucket.ID),
				SourceType:       types.SourceType(bucket.SourceType),
				CreateAt:         bucket.CreateAt,
				PaymentAddress:   bucket.PaymentAddress,
				PrimarySpAddress: bucket.PrimarySpAddress,
				ReadQuota:        bucket.ReadQuota,
				BillingInfo: types.BillingInfo{
					PriceTime:              0,
					TotalChargeSize:        0,
					SecondarySpObjectsSize: nil,
				},
			},
		})
	}
	resp = &metatypes.MetadataServiceGetUserBucketsResponse{Buckets: res}
	return resp, nil
}