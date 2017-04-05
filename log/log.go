package log

import (
	"fmt"
	"io"
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

// SetOutputFile set logging to a file instead of the old logger output
func SetOutputFile(filename string) (err error) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return fmt.Errorf("unable to open file %s: %v", filename, err)
	}
	Debugf("Logs will now be writted on %s, leaving stdout", filename)
	SetOutput(f)
	return nil
}

// SetOutput set logging to a new io.Writer
func SetOutput(output io.Writer) {
	logger = log.New(output, "", 0)
}

func formatLog(verbosity Verbose, callerSkip int, format string, args ...interface{}) (string, []interface{}) {
	date := time.Now().Format(time.RFC3339)
	caller := ""

	// we want to know which file/line throw this function call
	if callerSkip > 0 && (verbosity == DEBUG || verbosity == ERROR || verbosity == CRITICAL) {
		if pc, callerFile, callerLineNumber, ok := runtime.Caller(callerSkip); ok {
			details := runtime.FuncForPC(pc)
			callerName := details.Name()
			caller = fmt.Sprintf(" %s#%d %s", callerFile, callerLineNumber, callerName)
		} else {
			caller = " <unable to have caller infos>"
		}
	}
	prefix := fmt.Sprintf("%s [%s]%s", date, strings.ToUpper(VerboseReverseMapping[verbosity]), caller)
	args = append([]interface{}{prefix}, args...)
	return "%s " + format, args
}

// Logln is the base function used for logging,
// allowing custom prefix and caller skip,
// logs are not formatted, and a new line will be append
func Logln(verbose Verbose, callerSkipBeforeReachingHere int, args ...interface{}) {
	if verbose <= Verbosity {
		_, args = formatLog(verbose, callerSkipBeforeReachingHere+1, "", args...)
		logger.Println(args...)
	}
}

// Logf is the same function as Logln,
// the only difference is that the log can be formated
func Logf(verbose Verbose, callerSkipBeforeReachingHere int, format string, args ...interface{}) {
	if verbose <= Verbosity {
		format, args = formatLog(verbose, callerSkipBeforeReachingHere+1, format, args...)
		logger.Printf(format, args...)
	}
}

// Debugln log debug messages
func Debugln(args ...interface{}) {
	Logln(DEBUG, 2, args...)
}

// Debugf log debug messages
func Debugf(format string, args ...interface{}) {
	Logf(DEBUG, 2, format, args...)
}

// Requestln log request messages
func Requestln(args ...interface{}) {
	Logln(REQUEST, 2, args...)
}

// Requestf log request messages
func Requestf(format string, args ...interface{}) {
	Logf(REQUEST, 2, format, args...)
}

// Infoln log info messages
func Infoln(args ...interface{}) {
	Logln(INFO, 2, args...)
}

// Infof log info messages
func Infof(format string, args ...interface{}) {
	Logf(INFO, 2, format, args...)
}

// Warningln log warning messages
func Warningln(args ...interface{}) {
	Logln(WARNING, 2, args...)
}

// Warningf log warning messages
func Warningf(format string, args ...interface{}) {
	Logf(WARNING, 2, format, args...)
}

// Errorln log error messages
func Errorln(args ...interface{}) {
	Logln(ERROR, 2, args...)
}

// Errorf log error messages
func Errorf(format string, args ...interface{}) {
	Logf(ERROR, 2, format, args...)
}

// Criticalln log critical messages
func Criticalln(args ...interface{}) {
	Logln(CRITICAL, 2, args...)
}

// Criticalf log critical messages
func Criticalf(format string, args ...interface{}) {
	Logf(CRITICAL, 2, format, args...)
}
