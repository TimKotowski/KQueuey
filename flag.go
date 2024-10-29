package kqueuey

import (
	"flag"
)

type FlagOpts struct {
	ConfigPath string
	Logging    Logging
}

func (f *FlagOpts) RegisterFlags() {
	flag.StringVar(&f.ConfigPath, "config", "/usr/local/etc/", "set config file path to read from, default: /usr/local/etc/")
	flag.StringVar(&f.Logging.Level, "logLevel", "info", "provides log level: debug, info, trace, error, warn, fatal")
	flag.StringVar(&f.Logging.Format, "logFormat", "json", "provides what format logs should be presented as")
}

func (f *FlagOpts) Parse() {
	flag.Parse()
}
