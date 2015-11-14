package lxt

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/craigmj/commander"
)

// IsDir returns true if the given path is a directory,
// false otherwise
func IsDir(p string) (bool, error) {
	stat, err := os.Stat(p)
	if nil != err {
		return false, err
	}
	return stat.IsDir(), nil
}

// ContainerizePath takes a path that potentially starts
// with containerName:/path/ and turns it into a full
// path on the local machine.
func ContainerizePath(p string) (string, error) {
	i := strings.Index(p, ":")
	if -1 == i {
		return p, nil
	}
	name, path := p[0:i], p[i+1:]

	cont, err := GetDefinedContainer(name)
	if nil != err {
		return "", err
	}

	if !filepath.IsAbs(path) {
		path = filepath.Join("/home/ubuntu/", path)
	}
	return filepath.Join(cont.ConfigPath(), name, "rootfs", path), nil
}

func CopyFile(src, dest string) error {
	var err error
	src, err = ContainerizePath(src)
	if nil != err {
		return err
	}
	dest, err = ContainerizePath(dest)
	if nil != err {
		return err
	}

	isSrcDir, err := IsDir(src)
	if nil != err {
		return err
	}
	if isSrcDir {
		return fmt.Errorf("cp does not yet support copying directory")
	}
	isDestDir, err := IsDir(dest)
	if isDestDir {
		dest = filepath.Join(dest, filepath.Base(src))
	}
	in, err := os.Open(src)
	if nil != err {
		return err
	}
	defer in.Close()
	out, err := os.Create(dest)
	if nil != err {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func CopyCommand() *commander.Command {
	fs := flag.NewFlagSet("cp", flag.ExitOnError)

	return commander.NewCommand(
		"cp",
		"Copy a file into the container",
		fs,
		func(args []string) error {
			if 2 != len(args) {
				return fmt.Errorf("Require a source and destination parameter, but only %d args found", len(args))
			}
			return CopyFile(args[0], args[1])
		})
}
