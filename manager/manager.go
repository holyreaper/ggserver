package manager

//Manager 管理类
type Manager struct {
}

//Init init mng
func Init() {
	GMng = make(map[int]*Manager)
}

//GMng 管理类
var GMng map[int]*Manager
