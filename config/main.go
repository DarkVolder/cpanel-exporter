package config

import (
	"flag"
	"log"
)

type Config struct {
	Interval          int
	IntervalHeavy     int
	ListenPort        int
	ListenAddress     string
	Bandwidth         bool
	DomainsConfigured bool
	FtpAccounts       bool
	Sessions          bool
	License           bool
	Meta              bool
}

func New() (config Config) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.IntVar(&config.Interval, "interval", 60, "Check interval duration 60s by default")
	flag.IntVar(&config.IntervalHeavy, "interval_heavy", 1800, "Bandwidth and other heavy checks interval, 1800s (30min) by default")
	flag.IntVar(&config.ListenPort, "port", 59117, "Metrics Listen Port")
	flag.StringVar(&config.ListenAddress, "listenAddress", "0.0.0.0", "Metrics Listen Address")
	flag.BoolVar(&config.Bandwidth, "bandwidth", false, "Bandwidth by user")
	flag.BoolVar(&config.DomainsConfigured, "domainsConfigured", false, "Domains configured")
	flag.BoolVar(&config.FtpAccounts, "ftpAccounts", false, "Ftp accounts count")
	flag.BoolVar(&config.Sessions, "sessions", false, "Sessions web and email")
	flag.BoolVar(&config.License, "license", false, "License expire date and max users")
	flag.BoolVar(&config.Meta, "meta", false, "Meta version and release")
	flag.Parse()
	return
}
