package lblog

import (
	"github.com/holyreaper/ggserver/glog"
)

var globbyLog *glog.GLog

func init() {
	globbyLog = glog.NewGLog("lobby_")
}

//GetLobbyLogger GetLobbyLogger
func GetLobbyLogger() *glog.GLog {
	return globbyLog
}
