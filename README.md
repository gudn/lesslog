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

## Protocol
``` protocol-buffer
syntax = "proto3";

message Operation {
  uint64 sn = 1;
  bytes data = 2;
}

// Also you can provide X-Limit header
message FetchRequest {
  string log_name = 1;
  uint64 since_sn = 2;
}

message FetchResponse {
  repeated Operation operations = 1;
}

message PushRequest {
  string log_name = 1;
  uint64 last_sn = 2; // last_sn in current log state
  repeated Operation operations = 3; // sn in ops is ignored
}

message PushResponse {
  bool success = 1; // invalid last_sn or log isn't exist
  uint64 last_sn = 2;
}

// Create or get last sn
message CreateRequest {
  string log_name = 1;
}

service Lesslog {
  rpc Fetch(FetchRequest) returns (FetchResponse);
  rpc Watch(FetchRequest) returns (stream FetchResponse);
  rpc Push(PushRequest) returns (PushResponse);
  rpc Create(CreateRequest) returns (PushResponse);
}
```
