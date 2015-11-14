package lxt

import (
	"errors"
	"flag"

	"github.com/craigmj/commander"
)

// FixIP sets the IP of the named container to the
// given ip in the string, or to the container's current
// IP if the container is running.
func FixIP(containerName string, ip string) error {
	return errors.New("Sorry, FixIP isn't working yet!")
	return nil
}

func FixIPCommand() *commander.Command {
	fs := flag.NewFlagSet("fixip", flag.ExitOnError)
	n := fs.String("n", "", "Name of container")
	ip := fs.String("ip", "", "IP address desired")
	return commander.NewCommand(
		"fixip",
		"Sets a fixed IP for a named container",
		fs,
		func(args []string) error {
			if "" == *n {
				return errors.New("You must specify a container (-n)")
			}
			if "" == *ip {
				return errors.New("You need to specify the IP address (-ip)")
			}
			return FixIP(*n, *ip)
		})
}
