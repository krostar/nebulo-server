package validator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var (
	errStringEmpty    = errors.New("missing or empty")
	errStringBadParam = errors.New("bad check parameter")
)

// String check input to fulfil specified param
func String(v interface{}, param string) (err error) {

	var checksMapping = checkMap{
		"omitempty": checkDefinition{
			checkFct:  checkStringOmitEmpty,
			omitError: true,
		},
		"nonempty": checkDefinition{
			checkFct:  checkStringNonEmpty,
			omitError: false,
		},
		"length": checkDefinition{
			checkFct:  checkStringLength,
			omitError: false,
		},
		"custom": checkDefinition{
			checkFct:  checkStringCustom,
			omitError: false,
		},
	}

	return validate(v, param, checksMapping)
}

func checkStringOmitEmpty(str string, param string) (err error) {
	if str == "" {
		return errStringEmpty
	}
	return nil
}

func checkStringNonEmpty(str string, param string) (err error) {
	if len(str) == 0 {
		return errStringEmpty
	}

	return nil
}

func checkStringLength(str string, param string) (err error) {
	comparaisonType := map[string]func(real int, wanted int) error{
		"min": func(real int, wanted int) error {
			if real >= wanted {
				return nil
			}
			return fmt.Errorf("value should have at least %d  caracters but has only %d", wanted, real)
		},
		"equal": func(real int, wanted int) error {
			if real == wanted {
				return nil
			}
			return fmt.Errorf("value should have %d  caracters but has %d", wanted, real)
		},
		"max": func(real int, wanted int) error {
			if real <= wanted {
				return nil
			}
			return fmt.Errorf("value should have at most %d  caracters but has %d", wanted, real)
		},
	}

	options := strings.Split(param, "|")
	if len(options) != 2 {
		return errStringBadParam
	}
	comparaison := options[0]
	wanted, err := strconv.Atoi(options[1])
	if err != nil {
		return errStringBadParam
	}

	if cmp, ok := comparaisonType[comparaison]; ok {
		if err := cmp(len(str), wanted); err != nil {
			return err
		}
	} else {
		return errStringBadParam
	}
	return nil
}

func checkStringCustom(str string, param string) (err error) {
	// make message-generation easier
	caractersType := map[string][]*unicode.RangeTable{
		"upper case": {unicode.Upper, unicode.Title},
		"lower case": {unicode.Lower},
		"numeric":    {unicode.Number, unicode.Digit},
		"special":    {unicode.Space, unicode.Symbol, unicode.Punct, unicode.Mark},
	}

	// convert param to some usable options
	options := strings.Split(param, "|")
	if len(options) != 4 {
		return errStringBadParam
	}
	stringArrayToIntArray := func(options []string) (convertedOptions []int, err error) {
		convertedOptions = make([]int, 4)
		for index, opt := range options {
			conv, err := strconv.Atoi(opt)
			if err != nil {
				return nil, fmt.Errorf("atoi conversion failed: %v", err)
			}
			convertedOptions[index] = conv
		}
		return convertedOptions, nil
	}
	usableOptions, err := stringArrayToIntArray(options)
	if err != nil {
		return errStringBadParam
	}

	caractersTypeNumberRequired := map[string]int{
		"lower case": usableOptions[0],
		"upper case": usableOptions[1],
		"numeric":    usableOptions[2],
		"special":    usableOptions[3],
	}

	return stringCustomTest(str, caractersType, caractersTypeNumberRequired)
}

func stringCustomTest(str string, caractersType map[string][]*unicode.RangeTable, caractersTypeNumberRequired map[string]int) (err error) {
	caractersTypeNumber := make(map[string]int, len(caractersType))

	// make the tests
	for _, caracter := range str {
		for cType, cClass := range caractersType {
			if unicode.IsOneOf(cClass, caracter) {
				caractersTypeNumber[cType]++
			}
		}
	}

	// compare with the configuration defined above
	for caracterType, requiredNbr := range caractersTypeNumberRequired {
		var effectiveNbr = 0
		if effectiveNbrTest, ok := caractersTypeNumber[caracterType]; ok {
			effectiveNbr = effectiveNbrTest
		}
		if effectiveNbr < requiredNbr {
			if err != nil {
				err = fmt.Errorf("%s, %d %s (only %d detected)", err, requiredNbr, caracterType, effectiveNbr)
			} else {
				err = fmt.Errorf("value must have at least %d %s (only %d detected)", requiredNbr, caracterType, effectiveNbr)
			}
		}
	}

	return err
}
