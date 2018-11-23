package glog

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"runtime"

	"github.com/holyreaper/ggserver/def"
)

const (
	//DEBUG debug mode
	DEBUG = "[DEBUG]"

	//INFO INFO mode
	INFO = "[INFO]["

	//WARNING WARNING mode
	WARNING = "[WARNING]["

	//FATAL FATAL mode
	FATAL = "[FATAL]["
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
func InitLog(tp def.ServerType, id def.SID) (err error) {
	if def.ServerTypeNormal == tp {
		gLogger = NewGLog("lobby_" + strconv.Itoa(int(id)) + "_")
	} else if def.ServerTypeDB == tp {
		gLogger = NewGLog("db_" + strconv.Itoa(int(id)) + "_")
	} else if def.ServerTypeCenter == tp {
		gLogger = NewGLog("center_" + strconv.Itoa(int(id)) + "_")
	} else {
		fmt.Printf("InitLog fatal unknown servertype %d ", tp)
		gLogger = nil
		err = errors.New("invalid tp")
	}
	return err
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

		glog.Log = log.New(io.MultiWriter(open, os.Stdout), "LOG", log.Ldate|log.Ltime)
		glog.Openfile = open
	}
	return glog.Log
}

//LogInfo info mode
func LogInfo(format string, v ...interface{}) {
	file, line := getPath()
	logfile := fmt.Sprintf("][%s:%d]", file, line)
	format = logfile + format
	gLogger.getLogger()
	gLogger.Log.SetPrefix(INFO)
	gLogger.Log.Printf(format+"\r\n", v...)

}

//LogDebug DEBUG mode
func LogDebug(format string, v ...interface{}) {
	file, line := getPath()
	logfile := fmt.Sprintf("][%s:%d]", file, line)
	format = logfile + format
	gLogger.getLogger()
	gLogger.Log.SetPrefix(DEBUG)
	gLogger.Log.Printf(format+"\r\n", v...)
}

//LogWaring WARNING mode
func LogWaring(format string, v ...interface{}) {
	file, line := getPath()
	logfile := fmt.Sprintf("][%s:%d]", file, line)
	format = logfile + format
	gLogger.getLogger()
	gLogger.Log.SetPrefix(WARNING)
	gLogger.Log.Printf(format+"\r\n", v...)
}

//LogFatal FATAL mode
func LogFatal(format string, v ...interface{}) {
	file, line := getPath()
	logfile := fmt.Sprintf("][%s:%d]", file, line)
	format = logfile + format
	gLogger.getLogger()
	gLogger.Log.SetPrefix(FATAL)
	gLogger.Log.Printf(format+"\r\n", v...)
}
func getPath() (file string, line int) {
	var ok bool
	_, file, line, ok = runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	return
}
func getFileString() (file string) {
	file, line := getPath()
	file = fmt.Sprintf("[%s:%d]", file, line)
	return
}
