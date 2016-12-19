package lxt

import (
	"flag"
	"fmt"
	"os"

	"github.com/craigmj/commander"
)

func SetSudoer(n, u string) error {
	upath, err := ContainerizePath(fmt.Sprintf(`%s:/etc/sudoers.d/%s`, n, u))
	if nil != err {
		return err
	}
	out, err := os.OpenFile(upath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0440)
	if nil != err {
		return err
	}
	defer out.Close()
	fmt.Fprintf(out, "%s ALL=(ALL:ALL) NOPASSWD:ALL\n", u)
	return nil
}

func SudoerCommand() *commander.Command {
	fs := flag.NewFlagSet("sudoer", flag.ExitOnError)
	n := fs.String(`n`, ``, `name of container`)
	u := fs.String(`u`, `ubuntu`, `name of user`)
	return commander.NewCommand(`sudoer`, `Set the user as a nopasswd sudoer`,
		fs,
		func([]string) error {
			return SetSudoer(*n, *u)
		})
}
