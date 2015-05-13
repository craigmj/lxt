package lxt

import (
	"errors"
	"flag"
	"os"

	"gopkg.in/lxc/go-lxc.v2"

	"github.com/craigmj/commander"
)

func ExistsCommand() *commander.Command {
	fs := flag.NewFlagSet("exists", flag.ExitOnError)
	name := fs.String("n", "", "Name of the container")
	return commander.NewCommand(
		"exists", "Check whether a container exists - return 0 (exists) or 1 (not) to bash",
		fs,
		func([]string) error {
			if "" == *name {
				return errors.New("You need to specify a container name (-n)")
			}
			cont, err := lxc.NewContainer(*name, lxc.DefaultConfigPath())
			if nil != err {
				return err
			}
			if cont.Defined() {
				os.Exit(0)
			}
			os.Exit(1)
			return nil
		})
}
