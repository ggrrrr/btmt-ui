package app

import (
	"fmt"
)

type UnsupportedBodyTypeError struct {
	t string
}

func (e *UnsupportedBodyTypeError) Error() string {
	return fmt.Sprintf("unsupported body type %s", e.t)
}

type TemplateError struct {
	err error
}

func (e *TemplateError) Error() string {
	return fmt.Sprintf("cant process template %v", e.err)
}
