package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"

	"github.com/krostar/nebulo/log"
	"github.com/krostar/nebulo/returncode"
)

var (
	goodJWTSecretOption = "aaaaaAAAAA11111!!!!!0000000000000000000000"
	goodAESSecretOption = "aaaAAA111!!!00000000000000000000"
	goodEnvOption       = "prod"
	errLoadShouldThrow  = errors.New("load should have returned an error")
)

func init() {
	log.SetOutput(ioutil.Discard)
	defaultConfigurationFilePaths = []string{}

	/*
		default defer to use to reset everything:

		oldArgs := os.Args
		oldStdout := os.Stdout
		defer func() {
			os.Args = oldArgs
			os.Stdout = oldStdout
			Config = nil
			tester = false
			defaultConfigurationFilePaths = []string{}
		}()
	*/

}

func TestGoodOptions(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
		Config = nil
		if err := os.Remove("out.log"); err != nil {
			t.Fatal(err)
		}
	}()

	// test with good arguments
	os.Args = []string{oldArgs[0],
		"-e", goodEnvOption,
		"--jwt-secret=" + goodJWTSecretOption,
		"--aes-secret=" + goodAESSecretOption,
		"--address=127.0.0.1",
		"--port=10423",
		"--logging-file=out.log",
		"--verbose=debug",
	}
	if err := Load(); err != nil {
		t.Fatal("load should not have returned an error:", err)
	}
	if Config.Environment != goodEnvOption || Config.JWTSecret != goodJWTSecretOption || Config.AESSecret != goodAESSecretOption {
		t.Fatal("loaded configuration is bad")
	}

}

func TestBadBasicOptions(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
		Config = nil
	}()

	// test 1 with bad arguments: bad syntax
	os.Args = []string{oldArgs[0],
		"-e", "bad",
		"--jwt-secret=bad",
		"--aes-secret=bad",
		"--verbose=bad",
	}
	if err := Load(); err == nil {
		t.Fatal(errLoadShouldThrow)
	}
	Config = nil

	// test 2 with bad arguments: bad logging file
	os.Args = []string{oldArgs[0],
		"-e", goodEnvOption,
		"--jwt-secret=" + goodJWTSecretOption,
		"--aes-secret=" + goodAESSecretOption,
		"--logging-file=/",
	}
	if err := Load(); err == nil {
		t.Fatal(errLoadShouldThrow)
	}
}

func TestRequiredOptions(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
		Config = nil
	}()

	// test with missing required arguments environment
	os.Args = []string{oldArgs[0], "--jwt-secret=" + goodJWTSecretOption, "--aes-secret=" + goodAESSecretOption}
	if err := Load(); err == nil {
		t.Fatal(errLoadShouldThrow)
	}
}

func TestUselessOptions(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
		Config = nil
	}()
	// test with useless arguments add the end
	os.Args = []string{oldArgs[0],
		"-e", goodEnvOption,
		"--jwt-secret=" + goodJWTSecretOption,
		"--aes-secret=" + goodAESSecretOption,
		"--useless",
	}
	if err := Load(); err == nil {
		t.Fatal(errLoadShouldThrow)
	}
	Config = nil

	// test with useless arguments add the end
	os.Args = []string{oldArgs[0],
		"-e", goodEnvOption,
		"--jwt-secret=" + goodJWTSecretOption,
		"--aes-secret=" + goodAESSecretOption,
		"-u",
	}
	if err := Load(); err == nil {
		t.Fatal(errLoadShouldThrow)
	}
	Config = nil

	// test with useless arguments add the end
	os.Args = []string{oldArgs[0],
		"-e", goodEnvOption,
		"--jwt-secret=" + goodJWTSecretOption,
		"--aes-secret=" + goodAESSecretOption,
		"useless",
	}
	if err := Load(); err == nil {
		t.Fatal(errLoadShouldThrow)
	}
}

func TestGoodDefaultFileLoadingOptions(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
		Config = nil
		defaultConfigurationFilePaths = []string{}
	}()

	// get the path to this package
	_, filename, _, _ := runtime.Caller(0)
	filename = filepath.Dir(filename)

	// test file loading
	defaultConfigurationFilePaths = []string{filename + "/../config.sample.ini"}
	os.Args = []string{oldArgs[0], "-e", goodEnvOption, "--jwt-secret=" + goodJWTSecretOption, "--aes-secret=" + goodAESSecretOption}
	if err := Load(); err != nil {
		t.Fatal("load should not have returned an error:", err)
	}
}

func TestBadDefaultFileLoadingOptions(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
		Config = nil
		tester = false
		defaultConfigurationFilePaths = []string{}
	}()

	os.Args = []string{oldArgs[0],
		"-e", goodEnvOption,
		"--jwt-secret=" + goodJWTSecretOption,
		"--aes-secret=" + goodAESSecretOption,
	}

	// test file loading with failure: existing but unusable
	tester = true
	defaultConfigurationFilePaths = []string{"/"}
	if err := Load(); err == nil {
		t.Fatal(errLoadShouldThrow)
	}

	// reset to do the other test without interferences
	Config = nil
	tester = false

	// test file loading with failure: missing or unreadable file
	defaultConfigurationFilePaths = []string{"/inexisting"}
	if err := Load(); err != nil {
		t.Fatal(errLoadShouldThrow)
	}
}

func TestBadConfigFileOptions(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
		Config = nil
		tester = false
	}()

	// test with config file with failure
	tester = true
	os.Args = []string{oldArgs[0],
		"-e", goodEnvOption,
		"--jwt-secret=" + goodJWTSecretOption,
		"--aes-secret=" + goodAESSecretOption,
		"--config-file=/inexisting",
	}
	if err := Load(); err == nil {
		t.Fatal(errLoadShouldThrow)
	}
}

func TestOverrideOptions(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
		Config = nil
	}()

	// get the path to this package
	_, filename, _, _ := runtime.Caller(0)
	filename = filepath.Dir(filename)

	// test with config file arguments to overload previous config
	os.Args = []string{oldArgs[0],
		"-e", goodEnvOption,
		"--jwt-secret=" + goodJWTSecretOption,
		"--aes-secret=" + goodAESSecretOption,
		"--config-file=" + filename + "/../config.sample.ini",
	}
	if err := Load(); err != nil {
		t.Fatal("load should not have returned an error:", err)
	}
	if Config.Environment != "dev" {
		t.Fatal("the configuration file should have changed the env configuration")
	}
}

func TestHelp(t *testing.T) {
	oldArgs := os.Args
	oldStdout := os.Stdout
	defer func() {
		os.Args = oldArgs
		os.Stdout = oldStdout
		Config = nil
		tester = false
	}()

	args := []string{oldArgs[0], "-h"}
	// check return code
	if err := checkReturnCode(args, "TestHelp", returncode.HELP); err != nil {
		t.Fatal(err)
	}

	os.Args = args
	os.Stdout = nil
	tester = true
	if err := Load(); err != nil {
		t.Fatal(err)
	}
}

func TestConfigGen(t *testing.T) {
	oldArgs := os.Args
	oldStdout := os.Stdout
	defer func() {
		os.Args = oldArgs
		os.Stdout = oldStdout
		Config = nil
		tester = false
	}()

	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	// prepare to clean the tmp dir we are going to use
	defer func() {
		if err = os.RemoveAll(tmpdir); err != nil {
			t.Fatal(err)
		}
	}()

	args := []string{os.Args[0], "--config-gen=" + tmpdir + "/toto.ini"}
	// check return code
	if err = checkReturnCode(args, "TestConfigGen", returncode.CONFIGGEN); err != nil {
		t.Fatal(err)
	}

	// check for errors
	os.Args = args
	os.Stdout = nil
	tester = true
	if err := Load(); err != nil {
		t.Fatal(err)
	}
}

func checkReturnCode(args []string, testName string, exitCodeWanted int) (err error) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = args

	// used to check return status, in another thread
	if os.Getenv("GOTEST_CHECK_RETURN") == "1" {
		if err = Load(); err != nil {
			// we don't care, we are in another thread
			return err
		}
		Config = nil
		return
	}

	// start a new call to avoid test leave
	cmd := exec.Command(os.Args[0], "-test.run="+testName)
	cmd.Env = append(os.Environ(), "GOTEST_CHECK_RETURN=1")
	err = cmd.Run()
	var exitCode int

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			return err
		}
	} else {
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}

	if exitCode != exitCodeWanted {
		errExit := fmt.Sprintf("with config gen param, program is supposed to exit %d but has returned %d", exitCodeWanted, exitCode)
		return errors.New(errExit)
	}
	return nil
}
