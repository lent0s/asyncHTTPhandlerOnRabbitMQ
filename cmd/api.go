package cmd

import "github.com/lent0s/asyncHTTPhandlerOnRabbitMQ/config"

// Run asynchronous HTTP handler with RabbitMQ.
// For correctly work needs already running RabbitMQ-server, and config data, like:
//
// ##  [logger]
// logPathFolder:              ./logs
// logFileName:                log
// logMaxFileSizeKB:           1024
//
// ##  [rabbitMQ]
// rConnect:                   amqp://guest:guest@localhost:5672
//
// ##  [server]
// serverHost:                 127.0.0.1
// serverPort:                 9000
func RunService(LogMaxFileSizeKB uint16,
	LogPathFolder, LogFileName, RConnect, ServerHost, ServerPort string) {

	cfg := config.Config{
		LogPathFolder:    LogPathFolder,
		LogFileName:      LogFileName,
		LogMaxFileSizeKB: LogMaxFileSizeKB,
		RConnect:         RConnect,
		ServerHost:       ServerHost,
		ServerPort:       ServerPort,
	}
	local(cfg)
}
