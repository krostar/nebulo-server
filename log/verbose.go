package log

// Verbose is the type of all the available verbose
type Verbose int

// QUIET will hide every information
const QUIET = Verbose(0)

// CRITICAL is a verbosity level
const CRITICAL = Verbose(1)

// ERROR is a verbosity level
const ERROR = Verbose(2)

// WARNING is a verbosity level
const WARNING = Verbose(3)

// INFO is a verbosity level
const INFO = Verbose(4)

// REQUEST is a verbosity level
const REQUEST = Verbose(5)

// DEBUG is the maximum verbosity level
const DEBUG = Verbose(6)

// VerboseMapping is the translation between the verbose common name and the verbose type
var VerboseMapping = map[string]Verbose{
	"quiet":    QUIET,
	"critical": CRITICAL,
	"error":    ERROR,
	"warning":  WARNING,
	"info":     INFO,
	"request":  REQUEST,
	"debug":    DEBUG,
}

// VerboseReverseMapping is the translation between the verbosity type and the verbose common name
var VerboseReverseMapping = map[Verbose]string{
	QUIET:    "quiet",
	CRITICAL: "critical",
	ERROR:    "error",
	WARNING:  "warning",
	INFO:     "info",
	REQUEST:  "request",
	DEBUG:    "debug",
}
