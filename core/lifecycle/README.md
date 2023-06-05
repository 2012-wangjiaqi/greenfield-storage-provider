# Lifecycle

Lifecycle interface manages the lifecycle of a service and tracks its state changes. It also listens for signals from
the process to ensure a graceful shutdown.

## Concept

### Service Interface

Service is an interface for Lifecycle to manage. The component that plans to use Lifecycle needs to implement this interface.

```go
// Service is the interface for ServiceLifecycle to manage.
type Service interface {
	// Name defines the unique identifier of the service, which cannot be repeated 
	// globally.
	Name() string
	// Start the service, for resource application, start background coroutine and 
	// other startup operations.
	//
	// The Start method should be used in non-block way, for example, a blocked 
	// listening socket should open a goroutine separately internally.
	Start(ctx context.Context) error
	// Stop the service, close the goroutines inside the service, recycle resources, 
	// and ensure the graceful launch of the service.
	Stop(ctx context.Context) error
}
```

# Example
```go
    ctx := context.Background()
    svcLifecycle.RegisterServices(service...)
	// blocks the svcLifecycle for waiting signals to shut down the process
    svcLifecycle.Signals(syscall.SIGINT, syscall.SIGTERM ...).Init(ctx).StartServices(ctx).Wait(ctx)
```
