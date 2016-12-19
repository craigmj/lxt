package main

import (
	"log"

	"github.com/craigmj/commander"

	"lxt"
)

func main() {
	if err := commander.Execute(nil,
		lxt.AnsibleFacts,
		lxt.AutoStartCommand,
		lxt.CopyCommand,
		lxt.CPathCommand,
		lxt.ExistsCommand,
		lxt.ExposePortCommand,
		lxt.FixIPCommand,
		lxt.GetIPCommand,
		lxt.HostCommand,
		lxt.LnCommand,
		lxt.PortForwardCommand,
		lxt.ShellCommand,
		lxt.SSHKeyCommand,
		lxt.SSHKeyScanCommand,
		lxt.SudoerCommand,
		lxt.TeeCommand,
		lxt.WriteAnsibleFact,
	); nil != err {
		log.Fatal(err)
	}
}
