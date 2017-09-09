package power

import (
	"fmt"
	"strings"
)

// Default source configuration
var (
	DefaultCommunity = "public"
	DefaultPort      = 161
	DefaultRetries   = uint(1)
)

// Source describes a source of power statists data. It holds the necessary
// information to perform an SNMP query.
//
// Sources are queried to produce statistics, which are then fed to
// destinations.
type Source struct {
	Name      string
	Address   string
	Community string // SNMP community name
	Retries   uint
}

// String returns a string encoded representation of the power source.
func (s Source) String() (value string) {
	if s.Community == "" {
		value = s.Address
	} else {
		value = fmt.Sprintf("%s@%s", s.Community, s.Address)
	}

	if s.Name != "" {
		value = fmt.Sprintf("[%s]%s", s.Name, value)
	}

	return value
}

// ParseSource parses the given string as a power source description.
//
// The address is expected to be in one of these forms:
//
//   [name]community@host:port
//   [name]community@host
//   [name]host:port
//   [name]host
//   community@host:port
//   community@host
//   host:port
//   host
//
// [ups2s]tripplite@ups2s:161,tripplite@ups2n:161
func ParseSource(s string) (src Source, err error) {
	if s == "" {
		err = fmt.Errorf("empty source address")
		return
	}

	// Name
	if strings.HasPrefix(s, "[") && len(s) > 1 {
		if elements := strings.SplitN(s[1:], "]", 2); len(elements) == 2 {
			src.Name = elements[0]
			s = elements[1]
		}
	}

	// Community, Address
	if elements := strings.SplitN(s, "@", 2); len(elements) == 2 {
		src.Community = elements[0]
		src.Address = elements[1]
	} else {
		src.Community = DefaultCommunity
		src.Address = s
	}

	// Retries
	src.Retries = DefaultRetries

	// Validation
	if src.Address == "" {
		err = fmt.Errorf("no address specified for source \"%s\"", s)
		return
	}

	// Port
	if !strings.Contains(src.Address, ":") {
		src.Address = fmt.Sprintf("%s:%d", src.Address, DefaultPort)
	}

	return
}

// ParseSources takes the given set of strings and attempts to parse each one
// as a power source description.
func ParseSources(s []string) (sources []Source, err error) {
	for i, item := range s {
		source, pErr := ParseSource(item)
		if pErr != nil {
			return nil, fmt.Errorf("unable to parse source address %d: %s", i, pErr)
		}
		sources = append(sources, source)
	}
	return
}
