package kqueuey

import (
	"flag"
)

type FlagOpts struct {
	ConfigPath string
	Logging    Logging
}

func (f *FlagOpts) RegisterFlags() {
	flag.StringVar(&f.ConfigPath, "config", "/usr/local/etc/", "path to configuration file, default: /usr/local/etc/")
	flag.StringVar(&f.Logging.Level, "logLevel", "info", "provides log level: debug, info, warn, error, default: info")
	flag.StringVar(&f.Logging.Format, "logFormat", "json", "logger format options")
}

func (f *FlagOpts) Parse() {
	flag.Parse()
}
