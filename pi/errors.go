package pi

import (
	"errors"
	"fmt"

	"github.com/paultibbetts/mythicbeasts-client-go/internal/transport"
)

// ErrEmptyIdentifier is returned when an identifier is not used.
// Identifiers are required for all Pi resources.
var ErrEmptyIdentifier = errors.New("identifier is required")

// ErrNotFound indicates the API reported that the target resource does not exist.
var ErrNotFound = transport.ErrNotFound

// ErrIdentifierConflict indicates the requested resource identifier
// has already been used.
type ErrIdentifierConflict struct {
	Identifier string
}

func (e *ErrIdentifierConflict) Error() string {
	return fmt.Sprintf("identifier %q already in use", e.Identifier)
}
