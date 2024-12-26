package main

import (
	TMTServer "TMT_Attachment/server"
)

// "go run ."
func main() {
	// create server from constructor
	serv := TMTServer.MakeTMTServer(5, 5, 10)
	// toggle verbose logging of messaging stats
	serv.ReportMessagingDiagnostics()
	// begin simulator
	serv.Start()
}
