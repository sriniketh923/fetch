# Requirements
Go version >= 1.22.1

https://go.dev/doc/install

# Usage
```
➜  fetch ./fetch -h
usage: fetch [<flags>]

Flags:
  -h, --help                  Show context-sensitive help (also try --help-long and --help-man).
      --config="config.yaml"  Configuration file path.
      --log-level="info"      Set log level to debug, warn, info or error.
      --metrics=8080          Port to expose prometheus metrics.
```

# Instructions to run the application
Running the following command creates a `fetch` binary

```
➜  fetch make build
go build -o fetch main.go
```

Use the fetch binary to run the application. Run with `--log-level="debug"` mode for verbose logging. Example output:
```
➜  fetch ./fetch
2024-09-17 18:09:18.073363000 -400 EDT ./fetch: INFO:  fetch.com has 67% availability
2024-09-17 18:09:18.073454000 -400 EDT ./fetch: INFO:  www.fetchrewards.com has 100% availability
2024-09-17 18:09:32.894159000 -400 EDT ./fetch: INFO:  fetch.com has 67% availability
2024-09-17 18:09:32.894229000 -400 EDT ./fetch: INFO:  www.fetchrewards.com has 100% availability
2024-09-17 18:09:47.881534000 -400 EDT ./fetch: INFO:  fetch.com has 67% availability
2024-09-17 18:09:47.881630000 -400 EDT ./fetch: INFO:  www.fetchrewards.com has 100% availability
```

Run unit tests
```
go test ./...
?       fetch-interview [no test files]
?       fetch-interview/internal/endpoint       [no test files]
ok      fetch-interview/internal/config (cached)
```