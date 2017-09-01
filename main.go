package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/scjalliance/power/consolerecipient"
	"github.com/scjalliance/power/power"
	"github.com/scjalliance/power/stathatrecipient"
)

const (
	defaultSource     = "localhost"
	defaultRecipients = "console"
)

func main() {
	var (
		sourceStr    = os.Getenv("SOURCE")
		community    = os.Getenv("COMMUNITY")
		intervalStr  = os.Getenv("INTERVAL")
		recipientStr = os.Getenv("RECIPIENT")
		verbose      bool
	)

	if sourceStr == "" {
		sourceStr = defaultSource
	}
	if community == "" {
		community = power.DefaultCommunity
	}
	if recipientStr == "" {
		recipientStr = defaultRecipients
	}

	flag.StringVar(&sourceStr, "s", sourceStr, "comma separated list of power sources to query, in form [name]community@server:port")
	flag.StringVar(&community, "c", community, "default SNMP community for sources")
	flag.StringVar(&intervalStr, "n", intervalStr, "interval between executions, blank for single execution")
	flag.StringVar(&recipientStr, "r", recipientStr, "comma separated list of output recipients")
	flag.BoolVar(&verbose, "v", verbose, "show responses to unsupported queries")
	flag.Parse()

	power.DefaultCommunity = community
	power.RegisterRecipientType(stathatrecipient.Parse, "stathat")
	power.RegisterRecipientType(consolerecipient.Parse, "console")

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("No statistics specified.")
		os.Exit(2)
	}

	sources, err := power.ParseSources(strings.Split(sourceStr, ","))
	if err != nil {
		fmt.Printf("Source parsing error: %s\n", err)
		os.Exit(2)
	}
	if len(sources) == 0 {
		fmt.Printf("No sources specified\n")
		os.Exit(2)
	}

	recipients, err := power.ParseRecipients(recipientStr)
	if err != nil {
		fmt.Printf("Recipients parsing error: %s\n", err)
		os.Exit(2)
	}
	if len(recipients) == 0 {
		fmt.Printf("No recipients specified\n")
		os.Exit(2)
	}

	stats, err := power.ParseStatistics(args)
	if err != nil {
		fmt.Printf("Statistics parsing error: %s\n", err)
		os.Exit(2)
	}
	if len(stats) == 0 {
		fmt.Printf("No statistics specified\n")
		os.Exit(2)
	}

	execute(sources, stats, recipients, verbose)
}

func execute(sources []power.Source, stats []power.Statistic, recipients []power.Recipient, verbose bool) {
	for i, source := range sources {
		//fmt.Printf("Source %d (%s):\n", i, source)
		for _, r := range recipients {
			if handler, ok := r.(power.SourceHandler); ok {
				handler.SendSource(i, source)
			}
		}
		var values []power.Value
		values, err := power.Query(source, stats...)
		for _, r := range recipients {
			if err == nil {
				for _, v := range values {
					if verbose || !power.IsNotSupported(v.Err) {
						r.Send(v)
						//fmt.Printf("  %s: %s\n", v.Stat.Name, v)
					}
				}
			} else {
				if handler, ok := r.(power.ErrorHandler); ok {
					handler.SendQueryError(i, source, err)
				}
			}
		}
	}
}
