package power

import (
	"fmt"
	"strconv"
	"time"
)

// Value represents a single statistical value from a device or power source.
type Value struct {
	Source Source
	Stat   Statistic
	Time   time.Time
	Value  float64
	Err    error
}

// String returns a string representation of the value.
func (v Value) String() string {
	if v.Err != nil {
		return v.Err.Error()
	}
	return fmt.Sprintf("%s %s", strconv.FormatFloat(v.Value, 'f', -1, 64), v.Stat.Unit)
}

// StatName returns the name of the statistic.
func (v Value) StatName() string {
	sname := v.Source.Name
	if v.Source.Name == "" {
		sname = v.Source.Address
	}
	return fmt.Sprintf("%s %s", sname, v.Stat.Name)
}
