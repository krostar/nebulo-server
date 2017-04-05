package validator

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	errFilePathEmpty  = errors.New("filename is missing or empty")
	errFileUnwritable = errors.New("unable to write file")
)

// File check input to fulfil specified param
func File(v interface{}, param string) (err error) {
	var checksMapping = checkMap{
		"omitempty": checkDefinition{
			checkFct:  checkFileOmitEmpty,
			omitError: true,
		},
		"readable": checkDefinition{
			checkFct:  checkFileReadable,
			omitError: false,
		},
		"writable": checkDefinition{
			checkFct:  checkFileWritable,
			omitError: false,
		},
	}

	return validate(v, param, checksMapping)
}

func checkFileOmitEmpty(str string, param string) (err error) {
	if str == "" {
		return errFilePathEmpty
	}
	return nil
}

func checkFileReadable(str string, param string) (err error) {

	flag := os.O_RDONLY
	if param != "" {
		options := strings.Split(param, "|")
		for _, option := range options {
			switch option {
			case "createifmissing":
				flag |= os.O_CREATE
			}
		}
	}

	file, err := os.OpenFile(str, flag, 0600)
	if err != nil {
		return fmt.Errorf("unable to open file %s: %v", str, err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()
	return nil
}

func checkFileWritable(str string, param string) (err error) {
	file, err := os.OpenFile(str, os.O_RDWR, 0600)
	if err != nil {
		return fmt.Errorf("unable to open file %s: %s", str, err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	// save file state before doing something
	offset, readed, fileIsEmpty, err := saveFileStateBeforeWritting(file)
	if err != nil {
		return err
	}

	// try to write something
	if w, errW := file.Write([]byte("\x00")); w != 1 || errW != nil {
		return errFileUnwritable
	}

	// if file was empty, truncate(0) to avoid file length changement
	if fileIsEmpty {
		if err = os.Truncate(str, 0); err != nil {
			return fmt.Errorf("unable to truncate file %s: %v", str, err)
		}
	} else {
		if w, errW := file.WriteAt(readed, offset); w != 1 || errW != nil {
			return errFileUnwritable
		}
	}

	return nil
}

func saveFileStateBeforeWritting(file *os.File) (offset int64, readed []byte, fileIsEmpty bool, err error) {
	// make sure we know where to write to remove what we write later
	bof, err := file.Seek(0, 0)
	if err != nil {
		return 0, nil, false, fmt.Errorf("unable to seek: %v", err)
	}

	// we want to leave file as we found it
	fileIsEmpty = false
	readed = make([]byte, 1)

	if _, err = file.ReadAt(readed, 0); err != nil {
		if err == io.EOF {
			fileIsEmpty = true
			err = nil
		} else {
			return 0, nil, false, fmt.Errorf("unable to read: %v", err)
		}
	}
	return bof, readed, fileIsEmpty, err
}
