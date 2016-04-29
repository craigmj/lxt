package lxt

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	// "github.com/golang/glog"
	// "gopkg.in/lxc/go-lxc.v2"

	"github.com/craigmj/commander"
)

// addMountToLxc scans the lxc config for
// lxc.mount =
// It passes that through and a lxc.mount.entry as follows:
// lxc.mount.entry = /home/craig/proj/fundza/fundza var/proj/fundza none bind,optional,create=dir 0 0
// Once that's done, it continues, removing any existing mount for the named source
func addMountToLxc(C chan string, src, dest string) chan string {
	mountLine := fmt.Sprintf("lxc.mount.entry = %s %s none bind,optional,create=dir 0 0",
		src, dest)
	mountRegexp := regexp.MustCompile(`^\s*lxc\.mount\s*=`)
	duplicateRegexp := regexp.MustCompile(`^\s*lxc\.mount\.entry\s*=\s*` + src)
	return addLineToLxc(C, mountRegexp, duplicateRegexp, mountLine)
}

func LinkDirIntoContainer(n, src, dest string) error {
	// Check that src exists and convert to absolute path
	src, err := filepath.Abs(src)
	if nil != err {
		return err
	}
	// Stat the src and check that it is a directory
	fi, err := os.Stat(src)
	if nil != err {
		return err
	}
	if !fi.IsDir() {
		return errors.New("Your src path is not a directory")
	}

	cont, err := GetDefinedContainer(n)
	if nil != err {
		return err
	}

	isRunning := false
	if cont.Running() {
		isRunning = true
		if err = cont.Shutdown(10 * time.Second); nil != err {
			return fmt.Errorf("Failed to shutdown container: %s", err.Error())
		}
	}

	// Strip a leading path separator from the dest
	if os.PathSeparator == dest[0] {
		dest = dest[1:]
	}

	C, err := scanLxcConfig(cont.ConfigFileName())
	if nil != err {
		return err
	}
	if err = writeLxcConfig(cont.ConfigFileName(), addMountToLxc(C, src, dest)); nil != err {
		return err
	}

	if isRunning {
		// We've changed the configuration file, so we need to reload it
		if err = cont.LoadConfigFile(cont.ConfigFileName()); nil != err {
			return err
		}
		return cont.Start()
	}
	return nil
}

func LnCommand() *commander.Command {
	fs := flag.NewFlagSet("ln", flag.ExitOnError)
	n := fs.String("n", "", "Name of the container")
	src := fs.String("src", "", "Source directory to map into the container")
	dest := fs.String("dest", "", "Destination directory inside the container")

	return commander.NewCommand(
		"ln",
		"Link a host directory into a container",
		fs,
		func(args []string) error {
			var err error
			if "" == *n {
				return errors.New("You must name the container into which you want to link a directory (-n param)")
			}
			srcdir := *src
			if "" == srcdir {
				if 2 == len(args) {
					srcdir = args[0]
				} else {
					srcdir, err = os.Getwd()
					if nil != err {
						return errors.New("Cannot get cwd to use as source dir: " + err.Error())
					}
				}
			}
			destdir := *dest
			if "" == destdir {
				if 0 < len(args) {
					destdir = args[len(args)-1]
				} else {
					return errors.New("You must specify the destination directory inside the container (-dest)")
				}
			}
			return LinkDirIntoContainer(*n, srcdir, destdir)
		})
}
