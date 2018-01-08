package lblog

import (
	"github.com/holyreaper/ggserver/glog"
)

var globbyLog *glog.GLog

func init() {
	globbyLog = glog.NewGLog("lobby_")
}

//GetLogger GetLobbyLogger
func GetLogger() *glog.GLog {
	return globbyLog
}
