package envparse

import (
	"bufio"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	whiteSpaces = " \t"
)

// ParseEnv parses environment variables and sets properties on output
func ParseEnv(output interface{}) error {
	// start with Environment variables because they are lowest precedence
	for _, kv := range os.Environ() {
		for i, c := range kv {
			if c == '=' {
				if err := setConfig(output, kv[:i], kv[i+1:]); err != nil {
					return errors.Wrapf(err, "error with environment variable")
				}
				break
			}
		}
	}
	return nil
}

// ParseFile parses an Environment file and sets properties on output
func ParseFile(output interface{}, filePath string) error {
	// read configuration file
	file, err := os.Open(filePath)
	if err != nil {
		return errors.Wrapf(err, "could not read configuration file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNo := 1
	for scanner.Scan() {
		line := strings.TrimLeft(scanner.Text(), whiteSpaces)
		// skip commented lines and lines with 0 length
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		lineOk := false
	Line:
		for i, c := range line {
			if c == '=' {
				lineOk = true
				if err := setConfig(output, line[:i], line[i+1:]); err != nil {
					return errors.Wrapf(err, "error with configuration value")
				}
				break Line
			}
		}
		// lineOk is set to indicate that the line was in the correct format (i.e. key=value)
		if !lineOk {
			return errors.Errorf("Line %d isn't in the correct format. Expected 'KEY=value', Got '%s'", lineNo, line)
		}
		lineNo++
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func setConfig(conf interface{}, k string, v string) error {
	val := reflect.ValueOf(conf).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		if k == typeField.Tag.Get("config") {
			if !valueField.CanSet() {
				return errors.Errorf("Can't set value of %s because CanSet() returned false. See https://golang.org/pkg/reflect/#Value.CanSet for details", typeField.Name)
			}
			// valueField.Set(v)
			if valueField.Kind() == reflect.Int {
				i, err := strconv.Atoi(v)
				if err != nil {
					return errors.Wrapf(err, "%s expected int. Got '%s'", k, v)
				}
				x := int64(i)
				if valueField.OverflowInt(x) {
					return errors.Errorf("Number would overflow struct. %s", k)
				}

				valueField.SetInt(x)
			} else if valueField.Kind() == reflect.String {
				valueField.SetString(v)
			} else if valueField.Kind() == reflect.Bool {
				b, err := strconv.ParseBool(v)
				if err != nil {
					return errors.Wrapf(err, "%s expected bool. Got '%s'", k, v)
				}
				valueField.SetBool(b)
			} else {
				return errors.Errorf("Don't know how to set type '%s'", valueField.Kind())
			}
			break
		}
	}
	return nil
}
