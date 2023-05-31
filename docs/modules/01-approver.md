# Approver

Approver is responsible for processing the user's get create bucket and object approval 
request. 

## Interface 

```go
// Approver is the interface to handle ask approval.
type Approver interface {
	Modular
	// PreCreateBucketApproval prepares to handle CreateBucketApproval, it can do some
	// checks Example: check for duplicates, if limit specified by SP is reached, etc.
	PreCreateBucketApproval(ctx context.Context, task task.ApprovalCreateBucketTask) error
	// HandleCreateBucketApprovalTask handles the CreateBucketApproval, set expired
	// height and sign the MsgCreateBucket etc.
	HandleCreateBucketApprovalTask(ctx context.Context, task task.ApprovalCreateBucketTask) (bool, error)
	// PostCreateBucketApproval is called after HandleCreateBucketApprovalTask, it can
	// recycle resources, statistics and other operations.
	PostCreateBucketApproval(ctx context.Context, task task.ApprovalCreateBucketTask)

	// PreCreateObjectApproval prepares to handle CreateObjectApproval, it can do some
	// checks Example: check for duplicates, if limit specified by SP is reached, etc.
	PreCreateObjectApproval(ctx context.Context, task task.ApprovalCreateObjectTask) error
	// HandleCreateObjectApprovalTask handles the MsgCreateObject, set expired height
	// and sign the MsgCreateBucket etc.
	HandleCreateObjectApprovalTask(ctx context.Context, task task.ApprovalCreateObjectTask) (bool, error)
	// PostCreateObjectApproval is called after HandleCreateObjectApprovalTask, it can
	// recycle resources, statistics and other operations.
	PostCreateObjectApproval(ctx context.Context, task task.ApprovalCreateObjectTask)
	// QueryTasks queries tasks that running on approver by task sub key.
	QueryTasks(ctx context.Context, subKey task.TKey) ([]task.Task, error)
}
```

## Default Approver

### Get Create Bucket and Object Approval Workflow
* client make the `MsgCreateBucket` or `MsgCreateObject` and send the request to SP 
  with one of the above messages.
* SP check the `MsgCreateBucket` or `MsgCreateObject` message and set the expired height
  and signature that sign the message by SP operator private key.
* SP returns the `MsgCreateBucket` or `MsgCreateObject` message to client.
* client send the `MsgCreateBucket` or `MsgCreateObject` message to greenfield.
[TODO:: add create bucket and object docs of greenfield]().

### Get Create Bucket Approval
`PreCreateBucketApproval` will check the created buckets number by the same account, if
the number greater the config value(default 100) , or the get approval request is repeated, 
the request will be refused. 

`HandleCreateBucketApprovalTask` set the expired height(current block height adds timeout
height, the default timeout height is 10) and send the request to the signer module to sign,
then set signature that the result of signer return. Add the approval task in the task queue,
when the approval

`PostCreateBucketApproval` 