package envparse

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// ParseFile parses an Environment file and sets env variables
func ParseFile(filePath string) error {
	// read configuration file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return Parse(bufio.NewScanner(file))
}

func ParseString(str string) error {
	return Parse(bufio.NewScanner(strings.NewReader(str)))
}

func Parse(sc *bufio.Scanner) error {
	if sc == nil {
		return errors.New("scanner cannot be nil")
	}

	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimLeft(sc.Text(), " \t")
		// skip commented lines and lines with 0 length
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		sp := strings.Split(line, "=")
		// decline lines which are in the wrong format
		if len(sp) < 2 {
			return fmt.Errorf("Line %d isn't in the correct format. Expected 'KEY=value', Got '%s'", lineNo, line)
		}
		if err := os.Setenv(sp[0], strings.Join(sp[1:], "=")); err != nil {
			return fmt.Errorf("setting line %d of config file to env: %s", lineNo, err)
		}
	}

	return sc.Err()
}
