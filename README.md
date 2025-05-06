# Webhook Receiver

## Running the application using docker

First build the docker image using following command:

```bash
docker build -t webhook-receiver \
  --build-arg BATCH_SIZE=5 \
  --build-arg BATCH_INTERVAL_SEC=30 \
  --build-arg POST_ENDPOINT="https://eoa2dg5mkzbrlgw.m.pipedream.net" \
  .
```

Then run the docker container using following command:

```bash
docker run -d -p 8080:8080 webhook-receiver
```


## Testing the application

Testing health check endpoint:
```bash
curl http://localhost:8080/healthz
```

Testing log endpoint:

```bash
curl -X POST http://localhost:8080/log   -H "Content-Type: application/json"   -d '{
    "user_id": 1,
    "total": 1.65,
    "title": "delectus aut autem",
    "meta": {
      "logins": [
        {
          "time": "2020-08-08T01:52:50Z",
          "ip": "0.0.0.0"
        }
      ],
      "phone_numbers": {
        "home": "555-1212",
        "mobile": "123-5555"
      }
    },
    "completed": false
  }'
```


## Here are some of the test results:

```bash
docker logs inspiring_carver
2025/05/06 11:53:54 server.go:19: Initializing webhook receiver
2025/05/06 11:53:54 server.go:26: Configuration loaded: batch_size=5, batch_interval=30s, post_endpoint=https://eoa2dg5mkzbrlgw.m.pipedream.net
2025/05/06 11:53:54 server.go:43: Starting server on port 8080
2025/05/06 11:53:59 handlers.go:18: Request started: method=POST, path=/log, remote_addr=172.17.0.1
2025/05/06 11:53:59 handlers.go:23: Request completed: method=POST, path=/log, duration=570.789µs, status=202
2025/05/06 11:54:24 utils.go:104: Processing batch: size=1
2025/05/06 11:54:27 utils.go:135: Batch sent: size=1, status_code=200, duration=3.323507224s
```

## Results of golangci-lint

```bash
zahid@fedora:~/projects/Golang-Test-Project$ golangci-lint run
0 issues.
```





