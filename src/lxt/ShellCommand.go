package lxt

import (
	"errors"
	"flag"
	"io"
	"os"

	"github.com/craigmj/commander"
	"gopkg.in/lxc/go-lxc.v2"
)

func Shell(n, script string) error {
	c, err := GetDefinedContainer(n)
	if nil != err {
		return err
	}

	var sinFd uintptr
	if "-" == script {
		r, w, err := os.Pipe()
		if nil != err {
			return errors.New("Creating a pipe for stdin failed: " + err.Error())
		}
		go io.Copy(w, os.Stdin)
		sinFd = r.Fd()
	} else {
		in, err := os.Open(script)
		if nil != err {
			return errors.New("Failed to open " +
				script + " : " + err.Error())
		}
		defer in.Close()
		sinFd = in.Fd()
	}

	options := lxc.AttachOptions{
		UID:      os.Geteuid(),
		GID:      os.Getegid(),
		StdinFd:  sinFd,
		StdoutFd: os.Stdout.Fd(),
		StderrFd: os.Stderr.Fd(),
	}

	return c.AttachShell(options)
}

func ShellCommand() *commander.Command {
	fs := flag.NewFlagSet("shell", flag.ExitOnError)
	n := fs.String("n", "", "Name of container")
	return commander.NewCommand(
		"shell",
		"Run a local shell script inside the container: lxt shell -n CN test.sh",
		fs,
		func(args []string) error {
			for _, s := range args {
				if err := Shell(*n, s); nil != err {
					return err
				}
			}
			return nil
		})
}
