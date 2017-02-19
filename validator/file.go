package validator

import (
	"errors"
	"io"
	"os"
	"reflect"
	"strings"

	validator "gopkg.in/validator.v2"
)

type checkFunction func(filename string) (err error)

var checkMapping = map[string]checkFunction{
	"omitempty": checkOmitEmpty,
	"readable":  checkReadable,
	"writable":  checkWritable,
}

// File check input to respect param string
func File(v interface{}, param string) (err error) {
	// get the value to test
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return validator.ErrUnsupported
	}
	filename := st.String()

	checks := strings.Split(param, "+")
	for _, check := range checks {
		if checkFunc, ok := checkMapping[check]; ok {

			if err = checkFunc(filename); err != nil {
				if check == "omitempty" {
					return nil
				}
				return err
			}

		} else {
			return errors.New("unable to find file check named :'" + check + "'")
		}
	}

	return nil
}

func checkOmitEmpty(filename string) (err error) {
	if filename == "" {
		return errors.New("filename empty")
	}
	return nil
}

func checkReadable(filename string) (err error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	if err != nil {
		return errors.New(filename + ": " + err.Error())
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()
	return nil
}

func checkWritable(filename string) (err error) {
	file, err := os.OpenFile(filename, os.O_RDWR, 0600)
	if err != nil {
		return errors.New(filename + ": " + err.Error())
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	offset, readed, fileIsEmpty, err := saveFileStateBeforeWritting(file)
	if err != nil {
		return err
	}

	// try to write something
	if w, errW := file.Write([]byte("\x00")); w != 1 || errW != nil {
		return errors.New("unable to write")
	}

	// if file was empty, truncate(0) to avoid file length changement
	if fileIsEmpty {
		if err = os.Truncate(filename, 0); err != nil {
			return err
		}
	} else {
		if w, errW := file.WriteAt(readed, offset); w != 1 || errW != nil {
			return errors.New("unable to write")
		}
	}

	return nil
}

func saveFileStateBeforeWritting(file *os.File) (offset int64, readed []byte, fileIsEmpty bool, err error) {
	// make sure we know where to write to remove what we write later
	bof, err := file.Seek(0, 0)
	if err != nil {
		return 0, nil, false, err
	}

	// we want to leave file as we found it
	fileIsEmpty = false
	readed = make([]byte, 1, 1)

	if _, err = file.ReadAt(readed, 0); err != nil {
		if err == io.EOF {
			fileIsEmpty = true
		} else {
			return 0, nil, false, err
		}
	}
	return bof, readed, fileIsEmpty, err
}
