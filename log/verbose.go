package log

// Verbose is the type of all the available verbose
type Verbose int

const (
	// QUIET hide every information
	QUIET = Verbose(0)

	// CRITICAL hide most informations,
	// only the worst errors are shown
	CRITICAL = Verbose(1)

	// ERROR hide all non-essential informations,
	// only errors are shown
	ERROR = Verbose(2)

	// WARNING show all errors,
	// even thoses who aren't
	WARNING = Verbose(3)

	// INFO show every important steps of the workflow,
	// like server start, stop, update things, ...
	INFO = Verbose(4)

	// REQUEST show the same messages as INFO,
	// but show also the requests made by users
	REQUEST = Verbose(5)

	// DEBUG is the maximum verbosity level,
	// every pieces of useful informations are printed
	DEBUG = Verbose(6)
)

// VerboseMapping is the translation between
// the verbose common name and the verbose type
var VerboseMapping = map[string]Verbose{
	"quiet":    QUIET,
	"critical": CRITICAL,
	"error":    ERROR,
	"warning":  WARNING,
	"info":     INFO,
	"request":  REQUEST,
	"debug":    DEBUG,
}

// VerboseReverseMapping is the translation between
// the verbosity type and the verbose common name
var VerboseReverseMapping = map[Verbose]string{
	QUIET:    "quiet",
	CRITICAL: "critical",
	ERROR:    "error",
	WARNING:  "warning",
	INFO:     "info",
	REQUEST:  "request",
	DEBUG:    "debug",
}
