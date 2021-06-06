package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func runAction(path string, config actionConfig) error {
	cmd := exec.Command("node", fmt.Sprintf("%s/%s", path, config.Runs.Main))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	var errBuf strings.Builder
	cmd.Stderr = &errBuf

	err = cmd.Start()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	// TODO: are there multiline commands?
	for scanner.Scan() {
		if err := processOutputLine(scanner.Text()); err != nil {
			return err
		}
	}

	if err := cmd.Wait(); err != nil {
		if errStr := errBuf.String(); errStr != "" {
			return errors.Wrap(err, errStr)
		} else {
			return err
		}
	}
	return nil
}

func processOutputLine(line string) error {
	// TODO: support more commands:
	// set-env
	// add-mask
	// add-path
	switch {
	case !isCommand(line):
		fmt.Println(line)
		return nil
	case isDebug(line):
		return processDebug(line)
	case isWarning(line):
		return processWarning(line)
	case isError(line):
		return processError(line)
	case isSetOutput(line):
		return processSetOutput(line)
	default:
		fmt.Println(line)
		return nil
	}
}

var (

	// ::debug::debug level log
	// info level log
	// ::warning::warning level log
	// ::error::error level log
	// ::set-output name=greeting::Hello
	// ::set-env name=exportVarKey::exportVarValue
	// ::add-mask::secretToBeMasked
	// ::add-path::/path/to/mytool
	generalCommandRegexp *regexp.Regexp = regexp.MustCompile(`^::`)
	debugRegexp          *regexp.Regexp = regexp.MustCompile(`^::debug::`)
	warningRegexp        *regexp.Regexp = regexp.MustCompile(`^::warning::`)
	errorRegexp          *regexp.Regexp = regexp.MustCompile(`^::error::`)
	setOutputRegexp      *regexp.Regexp = regexp.MustCompile(`^::set-output name=(.*)::(.*)`)
)

func isCommand(line string) bool {
	return generalCommandRegexp.MatchString(line)
}

func isDebug(line string) bool {
	return debugRegexp.MatchString(line)
}

func isWarning(line string) bool {
	return warningRegexp.MatchString(line)
}

func isError(line string) bool {
	return errorRegexp.MatchString(line)
}

func isSetOutput(line string) bool {
	return setOutputRegexp.MatchString(line)
}

func processDebug(line string) error {
	// TODO: only print if debug mode set
	fmt.Println(debugRegexp.ReplaceAllString(line, "[DEBUG] "))
	return nil
}

func processWarning(line string) error {
	// TODO: use internal log lib for the CLI
	fmt.Println(warningRegexp.ReplaceAllString(line, "[WARNING] "))
	return nil
}

func processError(line string) error {
	// TODO: use internal log lib for the CLI
	fmt.Println(errorRegexp.ReplaceAllString(line, "[ERROR] "))
	return nil
}

func processSetOutput(line string) error {
	matches := setOutputRegexp.FindStringSubmatch(line)
	if len(matches) != 3 {
		return fmt.Errorf("set output command has unexpected format")
	}
	key := matches[1]
	value := matches[2] // TODO: should we uppercase this?
	cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", key, "--value", value).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set output environment variable, error: %#v | output: %s", err, cmdLog)
	}
	return nil
}
