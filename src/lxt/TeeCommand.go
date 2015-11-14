package lxt

import (
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/craigmj/commander"
)

func MakeDir(fn string) error {
	d := filepath.Dir(fn)
	return os.MkdirAll(d, 0755)
}

func MakeWriter(a bool, fn string) (out io.WriteCloser, err error) {
	fn, err = ContainerizePath(fn)
	if nil != err {
		return nil, err
	}
	if a {
		return os.OpenFile(fn, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	}
	err = MakeDir(fn)
	if nil != err {
		return
	}
	return os.Create(fn)
}

func Tee(a bool, out []string) (err error) {
	writers := make([]io.Writer, len(out), len(out)+1)
	for i, o := range out {
		writers[i], err = MakeWriter(a, o)
		if nil != err {
			return err
		}
	}
	writers = append(writers, os.Stdout)
	_, err = io.Copy(io.MultiWriter(writers...), os.Stdin)
	return err
}

func TeeCommand() *commander.Command {
	fs := flag.NewFlagSet("tee", flag.ExitOnError)
	a := fs.Bool("a", false, "Append to the named file.")
	return commander.NewCommand(
		"tee",
		"Just like Unix tee, but copies stdin into the destination in the container",
		fs,
		func(args []string) error {
			return Tee(*a, args)
		})
}
