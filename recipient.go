package power

import (
	"fmt"
	"strings"
)

var recipientTypes = make(map[string]RecipientParser)

// Recipient is something to which a power management value can be sent.
type Recipient interface {
	Send(Value)
}

// SourceHandler is a recipient that performs processing for each source.
type SourceHandler interface {
	SendSource(i int, s Source)
}

// ErrorHandler is a recipient that handles query errors.
type ErrorHandler interface {
	SendQueryError(i int, s Source, err error)
}

// RecipientParser is capable of parsing a given recipient address.
type RecipientParser func(address string) (Recipient, error)

// RegisterRecipientType registers the given recipient type for name.
//
// Name is case insensitive.
func RegisterRecipientType(parser RecipientParser, names ...string) {
	for _, name := range names {
		recipientTypes[strings.ToLower(name)] = parser
	}
}

// ParseRecipient parses the given string as a recipient definition in one of
// the following forms:
//
//   type
//   type:address
//
// The format of address depends on the value of type.
func ParseRecipient(s string) (recipient Recipient, err error) {
	if s == "" {
		err = fmt.Errorf("empty recipient definition")
		return
	}

	var (
		recipType string
		address   string
	)

	elements := strings.SplitN(s, ":", 2)
	if len(elements) == 2 {
		recipType, address = elements[0], elements[1]
	} else {
		recipType = s
	}

	parse, ok := recipientTypes[strings.ToLower(recipType)]
	if !ok {
		err = fmt.Errorf("unknown recipient type: \"%s\"", recipType)
		return
	}

	recipient, err = parse(address)
	return
}

// ParseRecipients takes the given set of strings and attempts to parse each one
// as a recipient.
func ParseRecipients(s string) (recipients []Recipient, err error) {
	elements := strings.Split(s, ",")
	for _, element := range elements {
		recipient, parseErr := ParseRecipient(element)
		if parseErr != nil {
			return nil, parseErr
		}
		recipients = append(recipients, recipient)
	}
	return
}
