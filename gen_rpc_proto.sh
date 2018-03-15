protoc.exe -I ./proto/ ./proto/ctrpc.proto --go_out=plugins=grpc:./ctrpc/
protoc.exe -I ./proto/ ./proto/dbrpc.proto --go_out=plugins=grpc:./dbrpc/
