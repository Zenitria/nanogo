package nanogo

import "fmt"

var (
	// ErrAccountNotFound ErrBlockNotFound is returned when the account isn't opened.
	ErrAccountNotFound = fmt.Errorf("account not found")
)
