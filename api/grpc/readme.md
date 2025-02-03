### gen golang code

From project root directory: 

```bash
$ protoc --go_out=internal/pkg/grpc/ --go_opt=paths=source_relative --go-grpc_out=internal/pkg/grpc/ --go-grpc_opt=paths=source_relative -I api/grpc channel_control.proto
```


### gen dart code

From project root directory:

```bash
$ protoc --dart_out=lib/src/services/proto -I server/api/grpc channel_control.proto
```