# Mfxkit - Mainflux Starter Kit

Mfxkit service provides a barebones HTTP API and Service interface implementation for development of a core Mainflux service.

## How-to

Copy `mfxkit` directory to the `mainflux` root directory, e.g. `~/go/src/github.com/mainflux/mainflux/`. Copy `cmd/mfxkit` directory to `mainflux/cmd` directory.

In `mainflux` root directory run

```bash
MF_MFXKIT_LOG_LEVEL=info go run cmd/mfxkit/main.go
```

You should get a message similar to this one

```bash
{"level":"info","message":"ping service http server listening at :9099 without TLS","ts":"2023-06-26T21:03:06.3116097Z"}
{"level":"info","message":"ping service gRPC server listening at :9199 without TLS","ts":"2023-06-26T21:03:06.311867558Z"}
```

In the other terminal window run

```bash
curl -i -X POST -H "Content-Type: application/json" localhost:9099/ping -d '{"secret":"secret"}'
```

If everything goes well, you should get

```bash
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 26 Jun 2023 21:03:26 GMT
Content-Length: 30

{"greeting":"Hello World :)"}
```

To change the secret or the port, prefix the `go run` command with environment variable assignments, e.g.

```bash
MF_MFXKIT_LOG_LEVEL=info MF_MFXKIT_SECRET=secret2 MF_MFXKIT_HTTP_PORT=9022 go run cmd/mfxkit/main.go
```

To see the change in action, run

```bash
curl -i -X POST -H "Content-Type: application/json" localhost:9022/ping -d '{"secret":"secret2"}'
```
