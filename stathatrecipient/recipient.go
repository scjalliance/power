package stathatrecipient

import (
	"fmt"

	"github.com/scjalliance/power"
)

// Recipient is a StatHat recipient of power management values.
type Recipient struct {
	ezkey string
}

// New returns a new StatHat recipient for the given ezkey.
func New(ezkey string) *Recipient {
	return &Recipient{
		ezkey: ezkey,
	}
}

// Send will send the report to the configured destination.
func (r *Recipient) Send(v power.Value) {
	if v.Err != nil {
		// Don't report stats that failed
		return
	}

	fmt.Printf("%s\n", v.StatName())
	//stathat.PostEZValueTime(v.StatName(), r.ezkey, v.Value, v.Time.Unix())
}

// Parse will parse the given address, which is a string containing the StatHat
// ezkey.
func Parse(address string) (power.Recipient, error) {
	return New(address), nil
}
