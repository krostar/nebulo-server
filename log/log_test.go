package log

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

var rexp *regexp.Regexp

type logXlnFct func(...interface{})
type logXfFct func(string, ...interface{})

var logMapLn map[Verbose]logXlnFct
var logMapF map[Verbose]logXfFct

type logRepr struct {
	Date    time.Time
	Verbose Verbose
	Caller  string
	Message string
}

func init() {
	var err error
	rexp, err = regexp.Compile(`^((?:(?:\d{4}-\d{2}-\d{2})T(?:\d{2}:\d{2}:\d{2}(?:\.\d+)?))(?:Z|[\+-]\d{2}:\d{2})?) \[(DEBUG|REQUEST|INFO|WARNING|ERROR|CRITICAL)\] (?:([^ ]+#\d+ [^[[:blank:]]+) )?(.*)\n$`)
	if err != nil {
		panic(err)
	}
	logMapLn = map[Verbose]logXlnFct{
		DEBUG:    Debugln,
		REQUEST:  Requestln,
		INFO:     Infoln,
		WARNING:  Warningln,
		ERROR:    Errorln,
		CRITICAL: Criticalln,
	}
	logMapF = map[Verbose]logXfFct{
		DEBUG:    Debugf,
		REQUEST:  Requestf,
		INFO:     Infof,
		WARNING:  Warningf,
		ERROR:    Errorf,
		CRITICAL: Criticalf,
	}
}

func TestFormat(t *testing.T) {
	var data string
	var err error
	var repr *logRepr

	// should have date prefix caller and message
	data = getStringFromLogOutput(DEBUG, "", "test1")
	repr, err = splitLogToRepr(data)
	if err != nil {
		t.Fatal(err)
	}
	if err = checkRepr(repr, DEBUG, false); err != nil {
		t.Fatal(err)
	}

	// should have date prefix and message
	data = getStringFromLogOutput(INFO, "test2: %d", 42)
	repr, err = splitLogToRepr(data)
	if err != nil {
		t.Fatal(err)
	}
	if err = checkRepr(repr, INFO, true); err != nil {
		t.Fatal(err)
	}
}

func TestAllLog(t *testing.T) {
	var data string
	var err error
	var repr *logRepr
	var emptyCaller bool

	for key := range VerboseReverseMapping {
		if key == QUIET {
			continue
		}
		if key == DEBUG || key == ERROR || key == CRITICAL {
			emptyCaller = false
		} else {
			emptyCaller = true
		}

		// logXln
		data = getStringFromLogOutput(key, "", "test1")
		repr, err = splitLogToRepr(data)
		if err != nil {
			t.Fatal(err, data)
		}
		if err = checkRepr(repr, key, emptyCaller); err != nil {
			t.Fatal(err, repr, data)
		}
		// logXf
		data = getStringFromLogOutput(key, "%s", "test2")
		repr, err = splitLogToRepr(data)
		if err != nil {
			t.Fatal(err, data)
		}
		if err = checkRepr(repr, key, emptyCaller); err != nil {
			t.Fatal(err, repr, data)
		}
	}
}

func TestOutputFile(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	// prepare to clean the tmp file we are going to use
	defer func() {
		SetOutput(os.Stdout)
		if err = os.Remove(tmpfile.Name()); err != nil {
			t.Fatal(err)
		}
	}()

	SetOutput(ioutil.Discard)

	// bad file
	if err = SetOutputFile("/missing"); err == nil {
		t.Fatal("SetOutputFile should throw")
	}

	// good file
	if err = SetOutputFile(tmpfile.Name()); err != nil {
		t.Fatal(err)
	}
	Infoln("test")

	data := make([]byte, 50)
	length, err := tmpfile.Read(data)
	if err != nil {
		t.Fatal(err)
	}
	if length <= 4 {
		t.Fatal("file length should be greater than 4 and is", length)
	}
}

func checkRepr(repr *logRepr, wantedVerbose Verbose, emptyCaller bool) (err error) {
	if repr.Verbose != wantedVerbose {
		return fmt.Errorf("verbose differs (want %d have %d)", wantedVerbose, repr.Verbose)
	}
	if len(repr.Caller) == 0 && !emptyCaller {
		return errors.New("caller should not be empty")
	} else if len(repr.Caller) != 0 && emptyCaller {
		return errors.New("caller should be empty")
	}
	return nil
}

func splitLogToRepr(log string) (repr *logRepr, err error) {
	repr = new(logRepr)

	parsed := rexp.FindStringSubmatch(log)
	if len(parsed) != 5 {
		return nil, errors.New("regex don't match")
	}

	date, err := time.Parse(time.RFC3339, parsed[1])
	if err != nil {
		return nil, err
	}
	repr.Date = date

	verbose, ok := VerboseMapping[strings.ToLower(parsed[2])]
	if !ok {
		return nil, errors.New("verbose isn't valid")
	}
	repr.Verbose = verbose

	repr.Caller = parsed[3]
	repr.Message = parsed[4]

	return repr, nil
}

func getStringFromLogOutput(verbosity Verbose, format string, args ...interface{}) (data string) {
	r, w := io.Pipe()
	SetOutput(w)
	defer func() {
		SetOutput(os.Stdout)
	}()

	go func() {
		defer func() {
			if err := w.Close(); err != nil {
				return
			}
		}()
		if format == "" {
			logMapLn[verbosity](args...)
		} else {
			logMapF[verbosity](format, args...)
		}
	}()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		return
	}
	data = buf.String()

	return data
}
