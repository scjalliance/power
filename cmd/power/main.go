package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/scjalliance/power"
	"github.com/scjalliance/power/consolerecipient"
	"github.com/scjalliance/power/stathatrecipient"
)

const (
	defaultStatistics = "all"
	defaultRecipients = "console"
)

func main() {
	var (
		statisticsStr = os.Getenv("STATISTICS")
		community     = os.Getenv("COMMUNITY")
		intervalStr   = os.Getenv("INTERVAL")
		recipientStr  = os.Getenv("RECIPIENT")
		interval      time.Duration
		verbose       bool
	)

	if statisticsStr == "" {
		statisticsStr = defaultStatistics
	}
	if community == "" {
		community = power.DefaultCommunity
	}
	if recipientStr == "" {
		recipientStr = defaultRecipients
	}

	//flag.StringVar(&sourceStr, "s", sourceStr, "comma separated list of power sources to query, in form [name]community@server:port")
	flag.StringVar(&statisticsStr, "q", statisticsStr, "comma separated list of statistics to query, \"all\" to include all statistics")
	flag.StringVar(&community, "c", community, "default SNMP community for sources")
	flag.StringVar(&intervalStr, "n", intervalStr, "interval between executions, blank for single execution")
	flag.StringVar(&recipientStr, "r", recipientStr, "comma separated list of output recipients")
	flag.BoolVar(&verbose, "v", verbose, "show responses to unsupported queries")
	flag.Parse()

	power.DefaultCommunity = community
	power.RegisterRecipientType(stathatrecipient.Parse, "stathat")
	power.RegisterRecipientType(consolerecipient.Parse, "console")

	sources, err := power.ParseSources(flag.Args())
	if err != nil {
		fmt.Printf("Source parsing error: %s\n", err)
		os.Exit(2)
	}
	if len(sources) == 0 {
		fmt.Printf("No sources specified\n")
		os.Exit(2)
	}

	recipients, err := power.ParseRecipients(strings.Split(recipientStr, ","))
	if err != nil {
		fmt.Printf("Recipients parsing error: %s\n", err)
		os.Exit(2)
	}
	if len(recipients) == 0 {
		fmt.Printf("No recipients specified\n")
		os.Exit(2)
	}

	stats, err := power.ParseStatistics(strings.Split(statisticsStr, ","))
	if err != nil {
		fmt.Printf("Statistics parsing error: %s\n", err)
		os.Exit(2)
	}
	if len(stats) == 0 {
		fmt.Printf("No statistics specified\n")
		os.Exit(2)
	}

	if intervalStr != "" {
		interval, err = time.ParseDuration(intervalStr)
		if err != nil {
			fmt.Printf("Unable to parse interval: %v", err)
			os.Exit(2)
		}
	}

	shutdown := NewShutdown()

	execute(sources, stats, recipients, verbose, shutdown)

	if interval > 0 {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				execute(sources, stats, recipients, verbose, shutdown)
			case <-shutdown:
				return
			}
		}
	}
}

func execute(sources []power.Source, stats []power.Statistic, recipients []power.Recipient, verbose bool, shutdown Shutdown) {
	for i, source := range sources {
		for _, r := range recipients {
			if shutdown.Signaled() {
				return
			}
			if handler, ok := r.(power.SourceHandler); ok {
				handler.SendSource(i, source)
			}
		}

		var values []power.Value
		values, err := power.Query(source, stats...)
		for _, r := range recipients {
			if err == nil {
				for _, v := range values {
					if shutdown.Signaled() {
						return
					}
					if verbose || !power.IsNotSupported(v.Err) {
						r.Send(v)
					}
				}
			} else {
				if shutdown.Signaled() {
					return
				}
				if handler, ok := r.(power.ErrorHandler); ok {
					handler.SendQueryError(i, source, err)
				}
			}
		}
	}
}
