package power

import (
	"fmt"
	"net"
	"strconv"
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
	Host      string
	Port      string
	Community string // SNMP community name
	Retries   uint
}

// HostPort returns the combination of "host:port". Its format matches that of
// net.JoinHostPort().
func (s Source) HostPort() string {
	return net.JoinHostPort(s.Host, s.Port)
}

// String returns a string encoded representation of the power source.
func (s Source) String() (value string) {
	value = s.HostPort()

	if s.Community != "" {
		value = fmt.Sprintf("%s@%s", s.Community, value)
	}

	if s.Name != "" {
		value = fmt.Sprintf("%s~%s", value, s.Name)
	}

	return value
}

// ParseSource parses the given string as a power source description.
//
// The address is expected to be in one of these forms:
//
//   community@host:port~name
//   community@host~name
//   host:port~name
//   host~name
//   community@host:port
//   community@host
//   host:port
//   host
//
func ParseSource(s string) (src Source, err error) {
	if s == "" {
		err = fmt.Errorf("empty source address")
		return
	}

	// Name
	if elements := strings.SplitN(s, "~", 2); len(elements) == 2 {
		src.Name = elements[1]
		s = elements[0]
	}

	// Community, Host, Port
	var hostport string
	if elements := strings.SplitN(s, "@", 2); len(elements) == 2 {
		src.Community = elements[0]
		hostport = elements[1]
	} else {
		src.Community = DefaultCommunity
		hostport = s
	}

	// Remember that ipv6 hosts can look like this: "[::]:port"
	hasPort := strings.LastIndex(hostport, ":") > strings.LastIndex(hostport, "]")
	if hasPort {
		src.Host, src.Port, err = net.SplitHostPort(hostport)
		if err != nil {
			err = fmt.Errorf("invalid host port for source \"%s\": %v", s, err)
			return
		}
	} else {
		src.Host = hostport
		src.Port = strconv.Itoa(DefaultPort)
	}

	// Retries
	src.Retries = DefaultRetries

	// Validation
	if src.Host == "" {
		err = fmt.Errorf("no host address specified for source \"%s\"", s)
		return
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
