package glog

import (
	"log"
	"os"
	"strconv"
	"time"
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

func (glog *GLog) getFileName() string {
	now := time.Now()
	fileName := glog.FileHead + strconv.Itoa(now.Year()) + strconv.Itoa(int(now.Month())) + strconv.Itoa(now.Day()) + strconv.Itoa(now.Hour())
	return fileName
}

func (glog *GLog) getLogger() *log.Logger {
	filename := glog.getFileName()
	if glog.Filename != filename {
		if glog.Openfile != nil {
			glog.Openfile.Close()
		}
		open, err := os.Create(filename)
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
func (glog *GLog) LogInfo(format string, v ...interface{}) {
	glog.Log.SetPrefix(INFO)
	glog.Log.Printf(format, v...)
}

//LogDebug DEBUG mode
func (glog *GLog) LogDebug(format string, v ...interface{}) {
	glog.Log.SetPrefix(DEBUG)
	glog.Log.Printf(format, v...)
}

//LogWaring WARNING mode
func (glog *GLog) LogWaring(format string, v ...interface{}) {
	glog.Log.SetPrefix(WARNING)
	glog.Log.Printf(format, v...)
}

//LogFatal FATAL mode
func (glog *GLog) LogFatal(format string, v ...interface{}) {
	glog.Log.SetPrefix(FATAL)
	glog.Log.Printf(format, v...)
}
