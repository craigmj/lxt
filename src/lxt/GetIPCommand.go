package lxt

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"gopkg.in/lxc/go-lxc.v2"

	"github.com/craigmj/commander"
)

// GetIP returns the IPv4 of the named container, or
// an error if the container does not exist.
func GetIP(containerName string) (string, error) {
	c, err := GetDefinedContainer(containerName)
	if !c.Running() {
		err := c.Start()
		if nil != err {
			return "", err
		}
	}
	c.Wait(lxc.RUNNING, time.Second*(time.Duration)(30))
	ips, err := c.IPv4Addresses()
	if nil != err {
		return "", err
	}
	if 0 == len(ips) {
		return "", errors.New("no v4 ips for container " + containerName)
	}
	return ips[0], nil
}

func GetIPCommand() *commander.Command {
	fs := flag.NewFlagSet("getip", flag.ExitOnError)
	n := fs.String("n", "", "Name of container")
	return commander.NewCommand(
		"getip",
		"Gets the current IP of the named container",
		fs,
		func(args []string) error {
			ip, err := GetIP(*n)
			if nil != err {
				return err
			}
			fmt.Println(ip)
			return nil
		})
}
