#import pbmysql "github.com/holyreaper/ggserver/rpcservice/pb/mysql"
protoc.exe -I ./ -I ./ ./ctrpc.proto --go_out=plugins=grpc:../ctrpc
protoc.exe -I ./ ./dbrpc.proto  --go_out=plugins=grpc:../dbrpc
protoc.exe -I ./ ./mysql.proto  --go_out=plugins=grpc:../mysql
protoc.exe -I ./ ./common.proto  --go_out=plugins=grpc:../common
