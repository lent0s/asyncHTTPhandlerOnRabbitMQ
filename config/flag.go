package config

import (
	"flag"
)

func ReadFlag() *Config {

	return readFlag()
}

func readFlag() *Config {

	CFG := getConfig()

	flag.StringVar(&CFG.LogPathFolder,
		"lpf", CFG.LogPathFolder, "path to folder for log files")

	flag.StringVar(&CFG.LogFileName,
		"lfn", CFG.LogFileName, "preferred log file name")

	lmfs := flag.Uint("lms", uint(CFG.LogMaxFileSizeKB), "max size for log file")
	CFG.LogMaxFileSizeKB = uint16(*lmfs)

	flag.StringVar(&CFG.RConnect,
		"rmq", CFG.RConnect, "path to RabbitMQ server")

	flag.IntVar(&CFG.RTimeout,
		"rto", CFG.RTimeout, "RabbitMQ server timeout")

	flag.StringVar(&CFG.ServerHost,
		"sip", CFG.ServerHost, "server start IP-address")

	flag.StringVar(&CFG.ServerPort,
		"spn", CFG.ServerPort, "server port number")

	flag.Parse()
	return CFG
}
