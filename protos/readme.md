To generate the Go-specific code from proto. While in the root directory run
```
protoc ./protos/<package>/<protofile>  --go_out=./protos_gen/<package> --proto_path=.  --go-grpc_out=./protos_gen/<package>
```