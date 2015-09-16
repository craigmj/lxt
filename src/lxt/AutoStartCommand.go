package lxt

import (
	"errors"
	"flag"
	"fmt"
	"regexp"

	"github.com/craigmj/commander"
)

func SetContainerAutoStart(n string, autostart bool) error {
	autoN := 1
	if !autostart {
		autoN = 0
	}
	find := regexp.MustCompile(`^\s*lxc\.start\.auto\s*=`)
	return EditLxcAddLine(n, find, find, fmt.Sprintf("lxc.start.auto = %d", autoN))
}

func AutoStartCommand() *commander.Command {
	fs := flag.NewFlagSet("autostart", flag.ExitOnError)
	n := fs.String("n", "", "Name of container to set autostart")
	off := fs.Bool("off", false, "Turn off autostart for the container")
	return commander.NewCommand("autostart",
		"Set a container to autostart",
		fs,
		func([]string) error {
			if "" == *n {
				return errors.New("You need to specify a container (-n)")
			}
			return SetContainerAutoStart(*n, !(*off))
		})
}
