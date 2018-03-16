package common

import (
	"sync"

	"github.com/holyreaper/ggserver/def"
)

var gserverID def.SID

var gserverType def.ServerType

var gone sync.Once

//GetServerID .
func GetServerID() def.SID {
	return gserverID
}

//GetServerType .
func GetServerType() def.ServerType {
	return gserverType
}

//SetServerID .
func SetServerID(serverid def.SID) {
	gone.Do(func() {
		gserverID = serverid
		gserverType = def.GetServerType(serverid)
	})
}
