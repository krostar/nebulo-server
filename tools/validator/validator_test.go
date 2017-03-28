package validator

import (
	"io/ioutil"
	"os"
	"testing"

	validator "gopkg.in/validator.v2"
)

func TestValidator(t *testing.T) {
	// wrong validators

	// empty string
	if err := String(nil, ""); err == nil {
		t.Fatal("validate isn't supposed to succed to cast an empty string")
	}

	// empty string
	if err := String("test", ""); err == nil {
		t.Fatal("validate isn't supposed to succed to parse that ''")
	}

	// inexisting test
	if err := String("test", "omit"); err == nil {
		t.Fatal("validate isn't supposed to succed to parse that 'omit'")
	}

	// bad syntax 1
	if err := String("test", "+++"); err == nil {
		t.Fatal("validate isn't supposed to succed to parse that '+++'")
	}

	// bad syntax 2
	if err := String("test", "omitempty+"); err == nil {
		t.Fatal("validate isn't supposed to succed to parse that 'omitempty+'")
	}
}

func TestErrorGeneration(t *testing.T) {
	var err error
	if err = String("test", "omitempty+"); err == nil {
		t.Fatal("validate isn't supposed to succed to parse that 'omitempty+'")
	}
	_ = HTTPErrors(err)

	type test struct {
		A string `validate:"string:nonempty"`
		B string `validate:"file:readable"`
	}

	te := test{
		A: "",
		B: "/missing",
	}
	if err = validator.Validate(te); err != nil {
		_ = HTTPErrors(err)
	} else {
		t.Fatal("validator should have failed to validate this configuration")
	}

}

func TestFileReal(t *testing.T) {
	existingRDWRableFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	existingRDWRableFilepath := existingRDWRableFile.Name()
	existingRDableFilepath := "/etc/hosts"

	// prepare to clean the tmp file we are going to use
	defer func() {
		if err = os.Remove(existingRDWRableFilepath); err != nil {
			t.Fatal(err)
		}
	}()

	// omitempty
	if err := File(existingRDWRableFilepath, "omitempty"); err != nil {
		t.Fatal(err)
	}

	// readable
	if err := File(existingRDableFilepath, "readable"); err != nil {
		t.Fatal(err)
	}

	// readable
	if err := File(existingRDableFilepath, "readable:createifmissing"); err != nil {
		t.Fatal(err)
	}

	// writable
	if err := File(existingRDWRableFilepath, "writable"); err != nil {
		t.Fatal(err)
	}
	if _, err := existingRDWRableFile.WriteString("test"); err != nil {
		t.Fatal(err)
	}
	if err := File(existingRDWRableFilepath, "writable"); err != nil {
		t.Fatal(err)
	}

}

func TestFileFake(t *testing.T) {
	missingFilepath := "/missing"
	emptyFilePath := ""

	// omitempty
	if err := File(emptyFilePath, "omitempty+readable+writable"); err != nil {
		t.Fatal(err)
	}

	// readble
	if err := File(missingFilepath, "readable"); err == nil {
		t.Fatal("missing file isn't supposed to be readable")
	}

	// writable
	if err := File(emptyFilePath, "writable"); err == nil {
		t.Fatal("empty file isn't supposed to be writable")
	}
}

func TestStringLength(t *testing.T) {
	// omitempty
	if err := String("non-empty", "omitempty"); err != nil {
		t.Fatal(err)
	}
	if err := String("", "omitempty"); err != nil {
		t.Fatal(err)
	}

	// nonempty
	if err := String("", "nonempty"); err == nil {
		t.Fatal("empty string should not pass nonempty test")
	}
	if err := String("non-empty", "nonempty"); err != nil {
		t.Fatal(err)
	}

	// custom
	if err := String("not-whats-required0", "custom:20|2|2|2"); err == nil {
		t.Fatal("string should not pass custom test")
	}
	if err := String("missing param", "custom:20|2"); err == nil {
		t.Fatal("string should not pass custom test")
	}
	if err := String("bad param", "custom:20|a"); err == nil {
		t.Fatal("string should not pass custom test")
	}
	if err := String("", "custom:20|2|2|2"); err == nil {
		t.Fatal("string should not pass custom test")
	}
	if err := String("aaAA11!!-toto", "custom:2|2|2|2"); err != nil {
		t.Fatal(err)
	}
}

func TestStringOthers(t *testing.T) {
	if err := String("a", "length|3"); err == nil {
		t.Fatal("string should not pass length test")
	}
	if err := String("bad param", "length"); err == nil {
		t.Fatal("string should not pass length test")
	}
	if err := String("bad param", "length|a"); err == nil {
		t.Fatal("string should not pass length test")
	}
	if err := String("a", "length:min|3"); err == nil {
		t.Fatal("string should not pass length test")
	}
	if err := String("a", "length:equal|3"); err == nil {
		t.Fatal("string should not pass length test")
	}
	if err := String("aaaa", "length:max|3"); err == nil {
		t.Fatal("string should not pass length test")
	}
	if err := String("aaa", "length:min|3"); err != nil {
		t.Fatal(err)
	}
	if err := String("aaa", "length:equal|3"); err != nil {
		t.Fatal(err)
	}
	if err := String("aaa", "length:max|3"); err != nil {
		t.Fatal(err)
	}
}
