run everything

```
docker compose watch
```

stop the runner, so no jobs get picked up

```
docker compose down runner
```

run only the server

```
docker compose watch server
```

enqueue a job

```
curl http://localhost:8080/enqueue\?name\=d
```

list jobs

```
curl http://localhost:8080/jobs
```

run the node web app

```
docker compose watch web
```

check the web app

```
watch -n 1 "curl http://localhost:3000"
```

If the postgres password changes and you can't connect for some non-obvious reason, try this:

```
docker volume prune -a
```
