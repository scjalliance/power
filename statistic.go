package power

import (
	"fmt"
	"strings"

	"github.com/k-sone/snmpgo"
	"github.com/scjalliance/power/snmpvar"
)

// Statistic describes a single power management statistic.
type Statistic struct {
	Name   string          // Name of statistic
	Unit   string          // Unit of measurement
	OID    snmpgo.Oids     // One or more possible OID values for this statistic
	Mapper snmpvar.Float64 // SNMP value mapper
}

// ParseStatistic parses a statistic in string format and returns the parsed
// value.
//
// The statistic string format is a comma-separated list of optional elements,
// with a well-known statistic key followed by a series of property/value pairs:
//
//   KEY,name:NAME,oid:OID,unit:UNIT
//
// The format is intended to meet three goals:
//
// 1. Easy use of well known statistics with preconfigured values
// 2. Overriding properties in preconfigured values with custom variations
// 3. Specification of custom, proprietary or otherwise unspported statistics
//
// Examples:
//
//   "EstimatedMinutesRemaining"
//   "name:WidgetDuration,oid:OID"
//   "EstimatedMinutesRemaining,name:ZomboMinutes,oid:OID,unit:Unit"
func ParseStatistic(s string) (stat Statistic, err error) {
	if s == "" {
		err = fmt.Errorf("empty statistic description")
		return
	}

	elements := strings.Split(s, ",")
	if len(elements) == 0 {
		err = fmt.Errorf("malformed statistic description: \"%s\"", s)
		return
	}

	for _, element := range elements {
		if parts := strings.SplitN(element, ":", 2); len(parts) == 2 {
			property := parts[0]
			value := parts[1]
			switch strings.ToLower(property) {
			case "name":
				stat.Name = value
			case "unit":
				stat.Unit = value
			case "oid":
				oid, oidErr := snmpgo.NewOid(value)
				if oidErr != nil {
					err = fmt.Errorf("unable to parse oid \"%s\" in statistic description: %s", s, oidErr)
					return
				}
				stat.OID = snmpgo.Oids{oid}
			}
		} else {
			if lookup, found := statMap[strings.ToLower(element)]; found {
				stat = lookup
			} else {
				err = fmt.Errorf("unknown statistic \"%s\"", element)
				return
			}
		}
	}

	if len(stat.OID) == 0 {
		err = fmt.Errorf("no object ID specified within \"%s\"", s)
		return
	}

	if stat.Mapper == nil {
		stat.Mapper = snmpvar.Ident
	}

	return
}

// ParseStatistics takes the given set of strings and attempts to parse each one
// as a power statistic.
func ParseStatistics(s []string) (stats []Statistic, err error) {
	for i, item := range s {
		switch strings.ToLower(item) {
		case "all":
			stats = append(stats, statList...)
		default:
			stat, pErr := ParseStatistic(item)
			if pErr != nil {
				return nil, fmt.Errorf("unable to parse statistic definition %d: %s", i, pErr)
			}
			stats = append(stats, stat)
		}
	}
	return
}
