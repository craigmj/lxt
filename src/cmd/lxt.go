package main

import (
	"log"

	"github.com/craigmj/commander"

	"lxt"
)

func main() {
	if err := commander.Execute(nil, lxt.LinkCommand); nil != err {
		log.Fatal(err)
	}
}
