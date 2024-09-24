package app

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

type testError struct {
	err error
	msg string
}

func (t testError) Error() string {
	return fmt.Sprintf("%s: %v", t.msg, t.err)
}

func (t testError) Unwrap() error {
	return t.err
}

func loadFile(n string) (string, error) {
	s, err := os.ReadFile(n)
	if err != nil {
		err = fmt.Errorf("shit %w", err)
		return "", fmt.Errorf("readFile[%s] %w", n, err)
	}
	return string(s), nil
}

func parseFile(n string) error {
	_, err := loadFile(n)
	if err != nil {
		return testError{
			err: err,
			msg: "parsing not possible",
		}
		// return fmt.Errorf("parseFile %w", err)
	}
	return nil
}

func TestError(t *testing.T) {

	// _, err := loadFile("asd")
	err := parseFile("Asd")
	fmt.Printf("%#v \n\n", err)
	fmt.Printf("%v \n\n", err)
	var osError *os.PathError
	fmt.Printf("%#v \n\n", errors.As(err, &osError))

}
