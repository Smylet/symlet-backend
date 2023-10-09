# Task processing with Asynq

## How to Add a New Task:

### 1. Define Your Payload:

This is the data you want to process asynchronously. 

`task_name.go`:
```go
type PayloadYourTaskName struct {
    Field1 string    `json:"field1"`
    Field2 int       `json:"field2"`
    // ... other fields
}
```

### 2. Distribute the Task:

Create a method similar to `DistributeTaskSendVerifyEmail`:

`task_name.go`:
```go
func (distributor *RedisTaskDistributor) DistributeYourTaskName(
	ctx context.Context,
	payload *PayloadYourTaskName,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask("task:your_task_name", jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}
```

### 3. Process the Task:

Define a method to handle the task processing:

```go
func (processor *RedisTaskProcessor) ProcessYourTaskName(ctx context.Context, task *asynq.Task) error {
	var payload PayloadYourTaskName
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	// ... Your task processing logic here ...

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Msg("processed task")
	return nil
}
```

### 4. Add Task to the Task Processor:

You will need to register the task with the `asynq.ServeMux`:

```go
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc("task:your_task_name", processor.ProcessYourTaskName)
	// ... other task handlers ...

	return processor.server.Start(mux)
}
```

### 5. Integration:

When you want to add this task in your main logic or elsewhere - Where you want the task to be processed asynchronously:

- Define your payload.
- Call the `DistributeYourTaskName` method with the appropriate context, payload, and options.
  
Example:

```go
// Define payload
payload := PayloadYourTaskName{
	Field1: "SampleData",
	Field2: 1234,
}

// Distribute task
opts := []asynq.Option{
	asynq.MaxRetry(10),
	asynq.ProcessIn(5 * time.Second),
	asynq.Queue(worker.QueueDefault),
}

task.DistributeYourTaskName(ctx, &payload, opts...)
```

### 6. Add the DistributeYourTaskName func definition to TaskDistributor interface

 `distributor.go`:

```go


type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error,
    DistributeYourTaskName(
        ctx context.Context,
        payload *PayloadYourTaskName,
        opts ...asynq.Option,
    ) error
}

```