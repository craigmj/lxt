package main

import (
	"log"

	"github.com/craigmj/commander"

	"lxt"
)

func main() {
	if err := commander.Execute(nil,
		lxt.LinkCommand,
		lxt.ExistsCommand,
		lxt.AutoStartCommand,
		lxt.FixIPCommand,
		lxt.CopyCommand,
	); nil != err {
		log.Fatal(err)
	}
}
