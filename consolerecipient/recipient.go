package consolerecipient

import (
	"fmt"

	"github.com/scjalliance/power/power"
)

type recipient struct{}

// Recipient is a recipient that will print output to the console
var Recipient recipient

// Parse ignores the given address and returns a console recipient.
func Parse(address string) (power.Recipient, error) {
	return Recipient, nil
}

func (r recipient) Send(v power.Value) {
	//fmt.Fprintf(r.w, r.format, v.Stat.Name, v)
	fmt.Printf("  %s: %s\n", v.Stat.Name, v)
}

func (r recipient) SendSource(i int, s power.Source) {
	fmt.Printf("Source %d (%s):\n", i, s)
}

func (r recipient) SendQueryError(i int, s power.Source, err error) {
	fmt.Printf("  Error: %v\n", err)
}
