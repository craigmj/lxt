package main

import (
	"log"

	"github.com/craigmj/commander"

	"lxt"
)

func main() {
	if err := commander.Execute(nil,
		lxt.LnCommand,
		lxt.ExistsCommand,
		lxt.AutoStartCommand,
		lxt.GetIPCommand,
		lxt.FixIPCommand,
		lxt.HostCommand,
		lxt.CopyCommand,
		lxt.PortForwardCommand,
		lxt.ShellCommand,
		lxt.CPathCommand,
		lxt.TeeCommand,
	); nil != err {
		log.Fatal(err)
	}
}
