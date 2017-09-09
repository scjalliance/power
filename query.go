package power

import (
	"errors"
	"fmt"

	"github.com/k-sone/snmpgo"
	"github.com/scjalliance/power/snmpvar"
)

// Query will attempt to retrieve the source's statistics via SNMP.
func Query(source Source, stats ...Statistic) (results []Value, err error) {
	snmp, err := snmpgo.NewSNMP(snmpgo.SNMPArguments{
		Version:   snmpgo.V2c,
		Address:   source.HostPort(),
		Retries:   source.Retries,
		Community: source.Community,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create snmpgo.SNMP object: %s", err)
	}

	if err = snmp.Open(); err != nil {
		return nil, fmt.Errorf("failed to open connection: %s", err)
	}
	defer snmp.Close()

	for _, stat := range stats {
		value := Value{
			Source: source,
			Stat:   stat,
		}
		value.Value, value.Err = query(snmp, stat.OID, stat.Mapper)
		results = append(results, value)
	}
	return
}

// query sends an SNMP v2c request and parses the value contained in the
// response.
func query(snmp *snmpgo.SNMP, oids snmpgo.Oids, mapper snmpvar.Float64) (value float64, err error) {
	// Execute the query
	pdu, err := snmp.GetRequest(oids)
	if err != nil {
		err = fmt.Errorf("failed to execute SNMP request: %s", err)
		return
	}
	if pdu.ErrorStatus() != snmpgo.NoError {
		err = fmt.Errorf("SNMP agent returned an error: [%d] %s", pdu.ErrorIndex(), pdu.ErrorStatus())
		return
	}

	// Retrieve the variables (one varbind per requested object identifier)
	bindings := pdu.VarBinds()

	// Map the variables to a single value
	return varToValue(oids, bindings, mapper)
}

// varToValue scans the returned set of variables in priority order and returns
// the first one that's valid.
//
// If none of the variables are valid it returns the last error.
func varToValue(oids snmpgo.Oids, bindings snmpgo.VarBinds, mapper snmpvar.Float64) (value float64, err error) {
	for _, oid := range oids {
		binding := bindings.MatchOid(oid)
		if binding != nil {
			v := binding.Variable
			switch v.Type() {
			case "NoSucheInstance":
				err = ErrNoSuchInstance
			case "NoSucheObject":
				err = ErrNoSuchObject
			default:
				value, err = mapper(v)
				if err == nil {
					// We found a valid value, return it
					return
				}
				err = fmt.Errorf("unable to parse returned value %s: %s", v.String(), err)
			}
		}
	}

	// No valid values were found; make sure we return an error
	if err == nil {
		err = errors.New("no value returned")
	}

	return
}
