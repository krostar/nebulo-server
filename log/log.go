package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	logger *log.Logger

	// Verbosity is the maximum level of message to print/write
	Verbosity Verbose
)

func init() {
	Verbosity = DEBUG
	logger = log.New(os.Stdout, "", 0)
}

// SetFile set a file to log into instead of the standart output
func SetFile(filename string) (err error) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	Debugln("Logs will now be writted on " + filename + ", leaving stdout.")
	logger = log.New(f, "", 0)
	return nil
}

// Debugln log debug messages
func Debugln(args ...interface{}) {
	if DEBUG <= Verbosity {
		println(DEBUG, args...)
	}
}

// Debugf log debug messages
func Debugf(format string, args ...interface{}) {
	if DEBUG <= Verbosity {
		printf(DEBUG, format, args...)
	}
}

// Requestln log request messages
func Requestln(args ...interface{}) {
	if REQUEST <= Verbosity {
		println(REQUEST, args...)
	}
}

// Requestf log request messages
func Requestf(format string, args ...interface{}) {
	if REQUEST <= Verbosity {
		printf(REQUEST, format, args...)
	}
}

// Infoln log info messages
func Infoln(args ...interface{}) {
	if INFO <= Verbosity {
		println(INFO, args...)
	}
}

// Infof log info messages
func Infof(format string, args ...interface{}) {
	if INFO <= Verbosity {
		printf(INFO, format, args...)
	}
}

// Warningln log warning messages
func Warningln(args ...interface{}) {
	if WARNING <= Verbosity {
		println(WARNING, args...)
	}
}

// Warningf log warning messages
func Warningf(format string, args ...interface{}) {
	if WARNING <= Verbosity {
		printf(WARNING, format, args...)
	}
}

// Errorln log error messages
func Errorln(args ...interface{}) {
	if ERROR <= Verbosity {
		println(ERROR, args...)
	}
}

// Errorf log error messages
func Errorf(format string, args ...interface{}) {
	if ERROR <= Verbosity {
		printf(ERROR, format, args...)
	}
}

// Criticalln log critical messages
func Criticalln(args ...interface{}) {
	if CRITICAL <= Verbosity {
		println(CRITICAL, args...)
	}
}

// Criticalf log critical messages
func Criticalf(format string, args ...interface{}) {
	if CRITICAL <= Verbosity {
		printf(CRITICAL, format, args...)
	}
}

func formatLog(verbosity Verbose, format string, args ...interface{}) (string, []interface{}) {
	date := time.Now().Format(time.RFC3339)
	caller := ""

	if verbosity == DEBUG || verbosity == ERROR || verbosity == CRITICAL {
		caller = " <unable to have caller infos>"
		if pc, _, _, ok := runtime.Caller(3); ok {
			details := runtime.FuncForPC(pc)
			callerName := details.Name()
			callerFile, callerLineNumber := details.FileLine(pc)
			caller = fmt.Sprintf(" %s#%d %s", callerFile, callerLineNumber, callerName)
		}
	}
	prefix := fmt.Sprintf("%s [%s]%s", date, strings.ToUpper(VerboseReverseMapping[verbosity]), caller)
	args = append([]interface{}{prefix}, args...)
	return "%s" + format, args
}

func println(verbosity Verbose, args ...interface{}) {
	_, args = formatLog(verbosity, "", args...)
	logger.Println(args...)
}

func printf(verbosity Verbose, format string, args ...interface{}) {
	format, args = formatLog(verbosity, format, args...)
	logger.Printf(format, args...)
}
