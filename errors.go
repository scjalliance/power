package power

import "errors"

// SNMP errors
var (
	ErrNoSuchInstance = errors.New("not supported (no such instance)")
	ErrNoSuchObject   = errors.New("not supported (no such object)")
)

// IsNotSupported returns true if the error indicates an unsupported SNMP
// object type.
func IsNotSupported(err error) bool {
	if err == ErrNoSuchInstance {
		return true
	}
	if err == ErrNoSuchObject {
		return true
	}
	return false
}
