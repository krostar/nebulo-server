package log

import (
	"log"
	"os"
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

// Criticalln log critical messages
func Criticalln(args ...interface{}) {
	if CRITICAL <= Verbosity {
		println("[CRITICAL]", args...)
	}
}

// Criticalf log critical messages
func Criticalf(format string, args ...interface{}) {
	if CRITICAL <= Verbosity {
		printf("[CRITICAL]", format, args...)
	}
}

// Errorln log error messages
func Errorln(args ...interface{}) {
	if ERROR <= Verbosity {
		println("[ERROR]", args...)
	}
}

// Errorf log error messages
func Errorf(format string, args ...interface{}) {
	if ERROR <= Verbosity {
		printf("[ERROR]", format, args...)
	}
}

// Warningln log warning messages
func Warningln(args ...interface{}) {
	if WARNING <= Verbosity {
		println("[WARN]", args...)
	}
}

// Warningf log warning messages
func Warningf(format string, args ...interface{}) {
	if WARNING <= Verbosity {
		printf("[WARN]", format, args...)
	}
}

// Infoln log info messages
func Infoln(args ...interface{}) {
	if INFO <= Verbosity {
		println("[INFO]", args...)
	}
}

// Infof log info messages
func Infof(format string, args ...interface{}) {
	if INFO <= Verbosity {
		printf("[INFO]", format, args...)
	}
}

// Requestln log request messages
func Requestln(args ...interface{}) {
	if REQUEST <= Verbosity {
		println("[REQUEST]", args...)
	}
}

// Requestf log request messages
func Requestf(format string, args ...interface{}) {
	if REQUEST <= Verbosity {
		printf("[REQUEST]", format, args...)
	}
}

// Debugln log debug messages
func Debugln(args ...interface{}) {
	if DEBUG <= Verbosity {
		println("[DEBUG]", args...)
	}
}

// Debugf log debug messages
func Debugf(format string, args ...interface{}) {
	if DEBUG <= Verbosity {
		printf("[DEBUG]", format, args...)
	}
}

func println(prefix string, args ...interface{}) {
	date := time.Now().Format(time.RFC3339)

	args = append([]interface{}{date, prefix}, args...)
	logger.Println(args...)
}

func printf(prefix string, format string, args ...interface{}) {
	date := time.Now().Format(time.RFC3339)

	args = append([]interface{}{date, prefix}, args...)
	format = "%s %s " + format
	logger.Printf(format, args...)
}
