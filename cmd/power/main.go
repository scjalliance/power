package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/gentlemanautomaton/signaler"
	"github.com/scjalliance/power"
	"github.com/scjalliance/power/consolerecipient"
	"github.com/scjalliance/power/stathatrecipient"
)

const (
	defaultSource     = "localhost"
	defaultStatistics = "all"
	defaultRecipients = "console"
)

func main() {
	shutdown := signaler.New().Capture(os.Interrupt, syscall.SIGTERM)
	defer shutdown.Wait()
	defer shutdown.Trigger()

	var (
		sourceStr     = os.Getenv("SOURCE")
		statisticsStr = os.Getenv("STATISTICS")
		community     = os.Getenv("COMMUNITY")
		intervalStr   = os.Getenv("INTERVAL")
		recipientStr  = os.Getenv("RECIPIENT")
		interval      time.Duration
		verbose       bool
	)

	if sourceStr == "" {
		sourceStr = defaultSource
	}
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

	var sources []power.Source
	var err error
	if flag.NArg() > 0 {
		sources, err = power.ParseSources(flag.Args())
	} else {
		sources, err = power.ParseSources(strings.Split(sourceStr, ","))
	}
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

	execute(shutdown.Signal, sources, stats, recipients, verbose)

	if interval > 0 {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				execute(shutdown.Signal, sources, stats, recipients, verbose)
			case <-shutdown.Signal:
				return
			}
		}
	}
}

func execute(shutdown signaler.Signal, sources []power.Source, stats []power.Statistic, recipients []power.Recipient, verbose bool) {
	if shutdown.Signaled() {
		return
	}

	stop := shutdown.Derive()
	defer stop.Wait()
	defer stop.Trigger()

	ctx := stop.Context()

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
		values, err := power.Query(ctx, source, stats...)
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
