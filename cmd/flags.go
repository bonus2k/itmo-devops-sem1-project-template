package main

import (
	"flag"
	"os"
)

var (
	srvAddr    string
	logLevel   string
	dbUsername string
	dbPassword string
	dbConn     string
)

func parseFlags() {
	flag.StringVar(&srvAddr, "a", "localhost:8080", "server address")
	flag.StringVar(&logLevel, "log", "info", "log level")
	flag.StringVar(&dbUsername, "u", "validator", "db username")
	flag.StringVar(&dbPassword, "p", "val1dat0r", "db password")
	flag.StringVar(&dbConn, "d", "localhost:5432/project-sem-1", "database host:port/db_name")

	flag.BoolFunc("h", "print this help", func(s string) error {
		flag.PrintDefaults()
		os.Exit(1)
		return nil
	})
}
