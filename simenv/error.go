package simenv

import "fmt"

// Error -
type Error struct {
	Context     string
	Description string
}

// Error -
func (err *Error) Error() string {
	return fmt.Sprintf("[SimEnv] %s: %s", err.Context, err.Description)
}
