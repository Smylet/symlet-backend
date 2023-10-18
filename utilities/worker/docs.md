# Integrating Asynchronous Tasks with Asynq in smylet-backend

This document guides you through the process of adding a new asynchronous task in the API. Each task comprises a payload definition, distribution logic, processing function, and server integration.

## Steps to Add a New Task

### 1. Define the Task Payload

Start by specifying the data structure that your task will process asynchronously.

In `task_name.go`:

```go
type PayloadYourTaskName struct {
    Field1 string `json:"field1"`
    Field2 int    `json:"field2"`
    // ... add other fields as needed
}
```

### 2. Distribute the Task

Set up the task distribution by creating a new method in `task_name.go`. This method will queue the task for processing.

```go
func (distributor *RedisTaskDistributor) DistributeYourTaskName(
    ctx context.Context,
    payload *PayloadYourTaskName,
    opts ...asynq.Option,
) error {
    // Convert the payload into JSON format
    jsonPayload, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal task payload: %w", err)
    }

    // Create a new task with the specified payload
    task := asynq.NewTask("task:your_task_name", jsonPayload, opts...)
    
    // Enqueue the task for processing
    info, err := distributor.client.EnqueueContext(ctx, task)
    if err != nil {
        return fmt.Errorf("failed to enqueue task: %w", err)
    }

    // Log task enqueuing details
    log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
        Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
    
    return nil
}
```

### 3. Process the Task

Implement the task processing logic. This function will execute the task when it's dequeued.

```go
func (processor *RedisTaskProcessor) ProcessYourTaskName(ctx context.Context, task *asynq.Task) error {
    // Parse the task payload
    var payload PayloadYourTaskName
    if err := json.Unmarshal(task.Payload(), &payload); err != nil {
        return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
    }

    // ... Implement your task processing logic here ...

    // Log task processing details
    log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
        Msg("processed task")
    
    return nil
}
```

### 4. Register the Task Handler

Integrate the task processing function into the Asynq server by registering it with the `asynq.ServeMux`.

```go
func (processor *RedisTaskProcessor) Start() error {
    mux := asynq.NewServeMux()

    // Register task processing handlers
    mux.HandleFunc("task:your_task_name", processor.ProcessYourTaskName)
    // ... other task handlers ...

    // Start the Asynq server with the registered handlers
    return processor.server.Start(mux)
}
```

### 5. Execute the Task

When you're ready to execute this task asynchronously within your application:

1. Define your task payload.
2. Call the `DistributeYourTaskName` method with the appropriate context, payload, and options.

```go
// Example to distribute your task

// Define task payload
payload := PayloadYourTaskName{
    Field1: "SampleData",
    Field2: 1234,
}

// Set task options
opts := []asynq.Option{
    asynq.MaxRetry(10),
    asynq.ProcessIn(5 * time.Second),
    asynq.Queue(worker.QueueDefault),
}

// Distribute task
task.DistributeYourTaskName(ctx, &payload, opts...)
```

### 6. Update Interfaces

Don't forget to update the `TaskDistributor` and `TaskProcessor` interfaces to include your new methods.

In `utilities/worker/interfaces.go`:

```go
type TaskDistributor interface {
    // ... other methods ...
    DistributeYourTaskName(ctx context.Context, payload *PayloadYourTaskName, opts ...asynq.Option) error
}

type TaskProcessor interface {
    // ... other methods ...
    ProcessYourTaskName(ctx context.Context, task *asynq.Task) error
}
```

---

Remember, it's important to maintain consistency in code documentation and commenting within your team to ensure that everyone can understand and follow the intended use and functionality. Also, ensure that any variable or file names used in the documentation match those in the actual codebase.