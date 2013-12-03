package log

import (
	"fmt"
	"io"
	golog "log"
	"os"
	"runtime"
	"strings"
	"time"
)

var fileHandle *os.File
var writer io.Writer
var maxLevel LogLevel

var StripPackage bool
var PadLocation bool
var HideLocation bool

const (
	Always = LogLevel(iota)
	Error
	Warn
	Info
	Debug
	Trace
)

type LogLevel uint8

func Configure(fileName string, appendFile bool) error {
	StripPackage = true
	PadLocation = false
	HideLocation = false

	if len(fileName) > 0 {
		fh, e := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0660)
		if e != nil {
			if appendFile {
				fh, e = os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0660)
			} else {
				e = os.Rename(fileName, fileName+time.Now().Format(".20060102-150405"))
				fh, e = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0660)
			}
		}
		if e != nil {
			return e
		}
		fileHandle = fh
		writer = io.MultiWriter(fh, os.Stdout)
	} else {
		writer = os.Stdout
	}
	golog.SetFlags(golog.Ltime)
	golog.SetOutput(writer)
	maxLevel = Info
	return nil
}

func SetMaxLevel(level LogLevel) {
	maxLevel = level
}

func SetMaxLevelString(level string) {
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
		levelString = "[E] "
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
	//golog.Printf(levelString+f.Name()+"(): "+format,v...)
	if !HideLocation {
		var where string
		if idx := strings.LastIndexAny(f.Name(), "/"); StripPackage && idx > 0 {
			where = f.Name()[idx+1:]
		} else {
			where = f.Name()
		}

		if PadLocation && len(where) < 32 {
			where = fmt.Sprintf("%-32s", where)
		}

		golog.Printf(levelString+"["+where+"] "+format, v...)
	} else {
		golog.Printf(levelString+format, v...)
	}
}
