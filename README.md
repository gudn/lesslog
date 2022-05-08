# lesslog

This is simple append-only log without any conflicts. Each client must know last
SN (sequence number) in log to write.

**Warning:** this service not check any security rules. Use proxy for control
this.

## Start
For development:

```sh
docker-compose up -d
make # apply migrations
make grpc-start
```

In production see `config.dev` file format and make your own config (now not
support NATS auth).
