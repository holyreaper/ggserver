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

//Login hello ...
func (loginServer *Login) Login(context.Context, *loginrpc.LoginRequest) (*loginrpc.LoginReply, error) {
	return &loginrpc.LoginReply{Message: "helo"}, nil
}
