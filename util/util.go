package util

import (
	. "github.com/holyreaper/ggserver/def"
)

//GetServerType ..
func GetServerType(id SID) ServerType {
	s := ServerType(id / SERVERBASEVALUE)
	return s
}
