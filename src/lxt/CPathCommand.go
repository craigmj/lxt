package lxt

import (
	"fmt"
	"github.com/craigmj/commander"
)

func CPathCommand() *commander.Command {
	return commander.NewCommand(
		"cpath",
		"Convert a container path to a local (host) path",
		nil,
		func(args []string) error {
			for _, a := range args {
				l, err := ContainerizePath(a)
				if nil != err {
					return err
				}
				fmt.Println(l)
			}
			return nil
		})
}
