package log

import (
	"fmt"
	"io"
	golog "log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var fileHandle *os.File
var writer io.Writer
var maxLevel LogLevel
var logFileName string

var StripPackage bool
var PadLocation bool
var HideLocation bool

var RollingFileSize int64 = 1048576
var CurrentFileSize int64
var RollingFileMux sync.Mutex
var RollingFileCount int = 10

const (
	Always = LogLevel(iota)
	Error
	Warn
	Info
	Debug
	Trace
)

type LogLevel uint8

func Configure(fileName string, rollingFileSize int64, rollingFileCount int) error {
	StripPackage = true
	PadLocation = false
	HideLocation = false

	RollingFileSize = rollingFileSize
	RollingFileCount = rollingFileCount
	
	logFileName = fileName

	if len(fileName) > 0 {

		fi, err := os.Stat(fileName)
		if err == nil {
			CurrentFileSize = fi.Size()
		}

		openFile(fileName)

	} else {
		writer = os.Stdout
		golog.SetFlags(golog.Ldate | golog.Ltime)
		golog.SetOutput(writer)
	}

	maxLevel = Info
	return nil
}

func openFile(fileName string) error {
	fh, e := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0660)
	if e != nil {
		fh, e = os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0660)
		if e != nil {
			return e
		}
	}
	fileHandle = fh
	writer = io.MultiWriter(fh, os.Stdout)
	golog.SetOutput(writer)
	return nil
}

func SetMaxLevel(level LogLevel) {
	maxLevel = level
}

func SetMaxLevelString(level string) {

	if len(level) == 0 {
		maxLevel = Info
		return
	}

	switch level[0] {
	case 'E', 'e':
		maxLevel = Error
		break
	case 'W', 'w':
		maxLevel = Warn
		break
	case 'I', 'i':
		maxLevel = Info
		break
	case 'D', 'd':
		maxLevel = Debug
		break
	case 'T', 't':
		maxLevel = Trace
		break
	}
}

func IsEnabled(level LogLevel) bool {
	if level > maxLevel {
		return false
	}
	return true
}

type LogMessage struct {
	Level LogLevel  `json:"level"`
	Ts    time.Time `json:"ts"`
	Where string    `json:"where"`
	Msg   string    `json:"msg"`
}

var OnLogChan chan *LogMessage

func Printf(level LogLevel, format string, v ...interface{}) {
	if level > maxLevel {
		return
	}

	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)

	levelString := "[A] "
	switch level {
	case Always:
		levelString = "[A] "
		break
	case Error:
		levelString = "[E] "
		break
	case Warn:
		levelString = "[W] "
		break
	case Info:
		levelString = "[I] "
		break
	case Debug:
		levelString = "[D] "
		break
	case Trace:
		levelString = "[T] "
		break
	}

	var finalMsg string

	//golog.Printf(levelString+f.Name()+"(): "+format,v...)
	if !HideLocation {
		var where string
		if idx := strings.LastIndexAny(f.Name(), "/"); StripPackage && idx > 0 {
			where = f.Name()[idx+1:]
		} else {
			where = f.Name()
		}

		if OnLogChan != nil {
			OnLogChan <- &LogMessage{level, time.Now(), where, fmt.Sprintf(format, v...)}
		}

		if PadLocation && len(where) < 32 {
			where = fmt.Sprintf("%-32s", where)
		}

		finalMsg = fmt.Sprintf(levelString+"["+where+"] "+format, v...)
	} else {

		if OnLogChan != nil {
			OnLogChan <- &LogMessage{level, time.Now(), "", fmt.Sprintf(format, v...)}
		}

		finalMsg = fmt.Sprintf(levelString+format, v...)
	}
	golog.Printf(finalMsg)

	if RollingFileSize > 0 {
		RollingFileMux.Lock()
		CurrentFileSize += int64(len(finalMsg))

		if CurrentFileSize > RollingFileSize {

			fileHandle.Close()
			rollLogs()

			openFile(logFileName)
			CurrentFileSize = 0
		}
		RollingFileMux.Unlock()
	}
}

func rollLogs() {
	for i := RollingFileCount; i > 0; i-- {
		if i == 1 {
			os.Rename(logFileName, logFileName+fmt.Sprintf(".%d", i))
		} else {
			os.Rename(logFileName+fmt.Sprintf(".%d", i-1), logFileName+fmt.Sprintf(".%d", i))
		}
	}
}
