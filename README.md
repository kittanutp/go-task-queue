# GO-TASK_QUEUE
    go run main.go

- GCP task queue implementation in GO
- The implementation using concurentcy on Queue Manager

# BODY Schema
    {
    "type": "POST",
    "url":"//URL_PATH//,
    "payload": {},
    "query": {},
    }

# TODO
- Graceful shutdown (wait for queue to empty before shutodown)
- Auto Scaling on QueueManager
- Unit Test
- DDOS Prevention
- Retry
