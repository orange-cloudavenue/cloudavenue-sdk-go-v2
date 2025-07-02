package cav

import "fmt"

func (e Endpoint) String() string {
	return fmt.Sprintf("[%s] %s %s %s %s",
		e.Category,
		e.Version,
		e.Name,
		e.Method,
		e.PathTemplate)
}

// String returns a string representation of the Endpoint.
func (e Category) String() string {
	return string(e)
}

// String returns a string representation of the Version.
func (e Version) String() string {
	return string(e)
}

// String returns a string representation of the Method.
func (e Method) String() string {
	return string(e)
}
