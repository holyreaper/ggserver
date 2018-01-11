package glog

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/holyreaper/ggserver/def"
)

const (
	//DEBUG debug mode
	DEBUG = "DEBUG"

	//INFO INFO mode
	INFO = "INFO"

	//WARNING WARNING mode
	WARNING = "WARNING"

	//FATAL FATAL mode
	FATAL = "FATAL"
)

//GLog ...
type GLog struct {
	Openfile *os.File

	Filename string
	FileHead string

	Log *log.Logger
}

// init ...
func init() {

}

var gLogger *GLog

//InitLog ...
func InitLog(tp def.ServerType, id def.SID) {
	if def.ServerTypeNormal == tp {
		gLogger = NewGLog("lobby_" + strconv.Itoa(int(id)))
	} else if def.ServerTypeDB == tp {
		gLogger = NewGLog("db_" + strconv.Itoa(int(id)))
	} else if def.ServerTypeCenter == tp {
		gLogger = NewGLog("center" + strconv.Itoa(int(id)))
	} else {
		fmt.Printf("InitLog fatal unknown servertype %d ", tp)
		gLogger = nil
	}
}

//NewGLog new glog
func NewGLog(fileHead string) *GLog {
	return &GLog{FileHead: fileHead}
}

func (glog *GLog) getFileName() string {
	now := time.Now()
	fileName := glog.FileHead + strconv.Itoa(now.Year()) + "_" + strconv.Itoa(int(now.Month())) + "_" + strconv.Itoa(now.Day()) + "_" + strconv.Itoa(now.Hour()) + ".log"
	return fileName
}

func (glog *GLog) getLogger() *log.Logger {
	filename := glog.getFileName()
	if glog.Filename != filename {
		if glog.Openfile != nil {
			glog.Openfile.Close()
		}
		glog.Filename = filename
		//open, err := os.Create(filename)
		open, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("getLogger filename %s open fail", filename)
			return nil
		}
		glog.Log = log.New(open, "rpcservice", log.Ldate|log.Ltime|log.Llongfile)
		glog.Openfile = open
	}
	return glog.Log
}

//LogInfo info mode
func LogInfo(format string, v ...interface{}) {
	gLogger.getLogger()
	gLogger.Log.SetPrefix(INFO)
	gLogger.Log.Printf(format+"\r\n", v...)
}

//LogDebug DEBUG mode
func LogDebug(format string, v ...interface{}) {
	gLogger.getLogger()
	gLogger.Log.SetPrefix(DEBUG)
	gLogger.Log.Printf(format+"\r\n", v...)
}

//LogWaring WARNING mode
func LogWaring(format string, v ...interface{}) {
	gLogger.getLogger()
	gLogger.Log.SetPrefix(WARNING)
	gLogger.Log.Printf(format+"\r\n", v...)
}

//LogFatal FATAL mode
func LogFatal(format string, v ...interface{}) {
	gLogger.getLogger()
	gLogger.Log.SetPrefix(FATAL)
	gLogger.Log.Printf(format+"\r\n", v...)
}
