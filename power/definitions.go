package power

import (
	"strings"

	"github.com/k-sone/snmpgo"
	"github.com/scjalliance/power/snmpvar"
)

// Output Source Enumeration
const (
	OutputSourceOther = iota + 1
	OutputSourceNone
	OutputSourceNormal
	OutputSourceBypass
	OutputSourceBattery
	OutputSourceBooster
	OutputSourceReducer
)

var (
	statMap  = make(map[string]Statistic)
	statList []Statistic
)

func registerStat(stat Statistic) {
	key := strings.ToLower(stat.Name)
	statMap[key] = stat
	statList = append(statList, stat)
}

func init() {
	registerStat(EstimatedMinutesRemaining)
	registerStat(EstimatedChargeRemaining)
	registerStat(BatteryVoltage)
	registerStat(BatteryTemperature)
	registerStat(OnBattery)
	registerStat(InputVoltage)
	registerStat(InputCurrent)
	registerStat(OutputVoltage)
	registerStat(OutputCurrent)
	registerStat(OutputPower)
	registerStat(OutputPercentLoad)
}

// Preconfigured power management statistics
var (
	EstimatedMinutesRemaining = Statistic{
		Name:   "EstimatedMinutesRemaining",
		Unit:   "minutes",
		OID:    snmpgo.Oids{snmpgo.MustNewOid("1.3.6.1.2.1.33.1.2.3.0")},
		Mapper: snmpvar.Ident,
	}
	EstimatedChargeRemaining = Statistic{
		Name:   "EstimatedChargeRemaining",
		Unit:   "%",
		OID:    snmpgo.Oids{snmpgo.MustNewOid("1.3.6.1.2.1.33.1.2.4.0")},
		Mapper: snmpvar.Ident,
	}
	BatteryVoltage = Statistic{
		Name:   "BatteryVoltage",
		Unit:   "volts (DC)",
		OID:    snmpgo.Oids{snmpgo.MustNewOid("1.3.6.1.2.1.33.1.2.5.0")},
		Mapper: snmpvar.Div(10),
	}
	BatteryTemperature = Statistic{
		Name:   "BatteryTemperature",
		Unit:   "°C",
		OID:    snmpgo.Oids{snmpgo.MustNewOid("1.3.6.1.2.1.33.1.2.7.0")},
		Mapper: snmpvar.Ident,
	}
	InputVoltage = Statistic{
		Name:   "InputVoltage",
		Unit:   "volts",
		OID:    snmpgo.Oids{snmpgo.MustNewOid("1.3.6.1.2.1.33.1.3.3.1.3.1")},
		Mapper: snmpvar.Ident,
	}
	InputCurrent = Statistic{
		Name:   "InputCurrent",
		Unit:   "volts",
		OID:    snmpgo.Oids{snmpgo.MustNewOid("1.3.6.1.2.1.33.1.3.3.1.4.1")},
		Mapper: snmpvar.Ident,
	}
	OnBattery = Statistic{
		Name:   "OnBattery",
		Unit:   "yes/no",
		OID:    snmpgo.Oids{snmpgo.MustNewOid("1.3.6.1.2.1.33.1.4.1.0")},
		Mapper: snmpvar.Match(OutputSourceBypass),
	}
	OutputVoltage = Statistic{
		Name:   "OutputVoltage",
		Unit:   "volts",
		OID:    snmpgo.Oids{snmpgo.MustNewOid("1.3.6.1.2.1.33.1.4.4.1.2.1")},
		Mapper: snmpvar.Ident,
	}
	OutputCurrent = Statistic{
		Name:   "OutputCurrent",
		Unit:   "amps",
		OID:    snmpgo.Oids{snmpgo.MustNewOid("1.3.6.1.2.1.33.1.4.4.1.3.1")},
		Mapper: snmpvar.Div(10),
	}
	OutputPower = Statistic{
		Name:   "OutputPower",
		Unit:   "watts",
		OID:    snmpgo.Oids{snmpgo.MustNewOid("1.3.6.1.2.1.33.1.4.4.1.4.1")},
		Mapper: snmpvar.Ident,
	}
	OutputPercentLoad = Statistic{
		Name:   "OutputPercentLoad",
		Unit:   "%",
		OID:    snmpgo.Oids{snmpgo.MustNewOid("1.3.6.1.2.1.33.1.4.4.1.5.1")},
		Mapper: snmpvar.Ident,
	}
)

// TODO: Add aliases to stats and then create an addMap function that uses them.
