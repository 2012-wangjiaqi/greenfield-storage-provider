# Modular

Modular is a complete logical module of SP. The SP framework is responsible
for the necessary interaction between modules. As for the implementation
of the module, it can be customized. Example, The SP framework stipulates
that ask object approval must be carried out before uploading an object,
whether agrees the approval, strategy can be customized that to agree or 
refuse approval by SP.

# Concept

## Front Modular
Front Modular handles the user's request, the gater will generate corresponding
task and send to Front Modular, the Front Modular need check the request
is correct. and after handle the task maybe some extra work is required.
So the Front Modular has three interfaces for each task type, `PreHandleXXXTask`,
`HandleXXXTask` and`PostHandleXXXTask`. Front Modular includes: `Authorizer`, 
`Approver`, `Downloader` and `Uploader`.

## Background Modular
Background Modular handles the SP inner task, since it is internally
generated, the correctness of the information can be guaranteed, so only
have one interface`HandleXXXTask`. Background Modular includes: `Manager`,
`TaskExecutor`, `P2P` and `Signer`.


# How to Customize Module

```go
// new your own CustomizedApprover instance that implement the Approver interface
//  NewCustomizedApprover must be func type: 
//      func(app *GfSpBaseApp, cfg *gfspconfig.GfSpConfig) (coremodule.Modular, error)
approver := NewCustomizedApprover(GfSpBaseApp, GfSpConfig)

// the Special Modular name is Predefined
gfspapp.RegisterModularInfo(model.ApprovalModularName, model.ApprovalModularDescription, approver)

// new GfSp framework app
gfsp, err := NewGfSpBaseApp(GfSpConfig, CustomizeApprover(approver))
if err != nil {
    return err
}

gfsp.Start(ctx)
// the GfSp framework will replace the default Approver with CustomizedApprover
```

