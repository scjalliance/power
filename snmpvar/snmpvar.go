// Package snmpvar maps SNMP variables to float64 values.
package snmpvar

import (
	"fmt"

	"github.com/k-sone/snmpgo"
)

// Float64 maps SNMP variables to float64 values.
type Float64 func(snmpgo.Variable) (float64, error)

// Ident is an identify function that returns SNMP values as a float64.
var Ident = func(v snmpgo.Variable) (float64, error) {
	if i, ok := v.(*snmpgo.Integer); ok {
		return float64(i.Value), nil
	}
	return 0, fmt.Errorf("unexpected non-integer SNMP variable type %s", v.Type())
}

// Mul multiplies SNMP values by a multiplier.
var Mul = func(multiplier float64) Float64 {
	return func(v snmpgo.Variable) (float64, error) {
		if i, ok := v.(*snmpgo.Integer); ok {
			return float64(i.Value) * multiplier, nil
		}
		return 0, fmt.Errorf("unexpected non-integer SNMP variable type %s", v.Type())
	}
}

// Div divides SNMP values by a divisor.
var Div = func(divisor float64) Float64 {
	return func(v snmpgo.Variable) (float64, error) {
		if i, ok := v.(*snmpgo.Integer); ok {
			return float64(i.Value) / divisor, nil
		}
		return 0, fmt.Errorf("unexpected non-integer SNMP variable type %s", v.Type())
	}
}

// Match returns 1 if an SNMP value matches an element of the given set.
var Match = func(set ...int) Float64 {
	m := make(map[int]struct{}, len(set))
	for _, element := range set {
		m[element] = struct{}{}
	}
	return func(v snmpgo.Variable) (float64, error) {
		if i, ok := v.(*snmpgo.Integer); ok {
			var value float64
			if _, found := m[int(i.Value)]; found {
				value = 1
			}
			return value, nil
		}
		return 0, fmt.Errorf("unexpected non-integer SNMP variable type %s", v.Type())
	}
}
