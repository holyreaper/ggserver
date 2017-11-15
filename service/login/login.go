package login

//import "context"
import "golang.org/x/net/context"
import "github.com/holyreaper/ggserver/pb/login"

func init() {
	// init
}

//Login 登录实现
type Login struct {
}

//SayHello hello ...
func (loginServer *Login) SayHello(context.Context, *loginrpc.HelloRequest) (*loginrpc.HelloReply, error) {
	return &loginrpc.HelloReply{Message: "helo"}, nil
}
