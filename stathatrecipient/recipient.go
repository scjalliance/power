package stathatrecipient

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/gentlemanautomaton/stathat"
	"github.com/scjalliance/power"
)

// DefaultFormat is the default StatHat statistic naming format.
const DefaultFormat = "{{.Source.Host}} {{.Stat.Name}}"

// Recipient is a StatHat recipient of power management values. It contains the
// ezkey and naming template.
type Recipient struct {
	ezkey string
	t     *template.Template // StatHat stat naming template parsed by text/template
}

// New returns a new StatHat recipient for the given ezkey. The default naming
// template is used.
func New(ezkey string) *Recipient {
	return &Recipient{
		ezkey: ezkey,
		t:     template.Must(template.New("stathat").Parse(DefaultFormat)),
	}
}

// NewWithNameTemplate returns a new StatHat recipient for the given ezkey and
// stat name template.
func NewWithNameTemplate(ezkey, nameTemplate string) (*Recipient, error) {
	t, err := template.New("stathat").Parse(nameTemplate)
	if err != nil {
		return nil, err
	}
	return &Recipient{
		ezkey: ezkey,
		t:     t,
	}, nil
}

// Send will send the report to the configured destination.
func (r *Recipient) Send(v power.Value) {
	if v.Err != nil {
		// Don't report stats that weren't collected successfully
		return
	}

	name := r.StatName(v)

	fmt.Printf("Sending data for \"%s\" to StatHat...", name)
	reporter := stathat.New().EZKey(r.ezkey)
	err := reporter.PostEZ(name, stathat.KindValue, v.Value, &v.Time)
	if err != nil {
		// TODO: Try again later?
		fmt.Printf("failed: %v\n", err)
	} else {
		fmt.Printf("done\n")
	}
}

// StatName returns the formatted name of the statistic in StatHat.
func (r *Recipient) StatName(v power.Value) string {
	var buf bytes.Buffer
	r.t.Execute(&buf, v)
	return buf.String()
}

// Parse will parse the given address, which is a string containing the StatHat
// ezkey.
func Parse(address string) (power.Recipient, error) {
	if elements := strings.SplitN(address, "~", 2); len(elements) == 2 {
		ezkey, tmpl := elements[0], elements[1]
		return NewWithNameTemplate(ezkey, tmpl)
	}
	return New(address), nil
}
