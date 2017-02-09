package log

// Verbose is the type of all the available verbose
type Verbose int

// CRITICAL is the minimum verbosity level
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
	"critical": CRITICAL,
	"error":    ERROR,
	"warning":  WARNING,
	"info":     INFO,
	"request":  REQUEST,
	"debug":    DEBUG,
}
