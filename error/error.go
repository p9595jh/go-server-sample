package server

import "fmt"

type NameNotFoundError struct {
	Name string
}

func (e *NameNotFoundError) Error() string {
	return fmt.Sprintf("book name '%s' not found", e.Name)
}
