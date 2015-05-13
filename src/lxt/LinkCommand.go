package lxt

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"gopkg.in/lxc/go-lxc.v2"

	"github.com/craigmj/commander"
)

func scanLxcConfig(filename string) (chan string, error) {
	inf, err := os.Open(filename)
	if nil != err {
		return nil, err
	}
	in := bufio.NewScanner(inf)
	C := make(chan string)
	go func() {
		defer inf.Close()
		defer close(C)
		for in.Scan() {
			C <- in.Text()
		}
	}()
	return C, nil
}

// addMountToLxc scans the lxc config for
// lxc.mount =
// It passes that through and a lxc.mount.entry as follows:
// lxc.mount.entry = /home/craig/proj/fundza/fundza var/proj/fundza none bind,optional,create=dir 0 0
// Once that's done, it continues, removing any existing mount for the named source
func addMountToLxc(C chan string, src, dest string) chan string {
	D := make(chan string)
	go func() {
		defer close(D)
		mountLine := fmt.Sprintf("lxc.mount.entry = %s %s none bind,optional,create=dir 0 0",
			src, dest)
		foundMount := false
		mountRegexp := regexp.MustCompile(`^\s*lxc\.mount\s*=`)
		duplicateRegexp := regexp.MustCompile(`^\s*lxc\.mount\.entry\s*=\s*` + src)
		for l := range C {
			if duplicateRegexp.MatchString(l) {
				if !foundMount {
					D <- mountLine
					foundMount = true
				}
			} else if mountRegexp.MatchString(l) {
				D <- l
				if !foundMount {
					D <- mountLine
					foundMount = true
				}
			} else {
				D <- l
			}
		}
		if !foundMount {
			D <- mountLine
		}
	}()
	return D
}

func writeLxcConfig(filename string, C chan string) error {
	var b bytes.Buffer
	for l := range C {
		fmt.Fprintln(&b, l)
	}
	if err := ioutil.WriteFile(filename, b.Bytes(), 0644); nil != err {
		return err
	}
	return nil
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

	cont, err := lxc.NewContainer(n, lxc.DefaultConfigPath())
	if nil != err {
		return err
	}
	if !cont.Defined() {
		return errors.New("No container named " + n + " is defined")
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

func LinkCommand() *commander.Command {
	fs := flag.NewFlagSet("link", flag.ExitOnError)
	n := fs.String("n", "", "Name of the container")
	src := fs.String("src", "", "Source directory to map into the container")
	dest := fs.String("dest", "", "Destination directory inside the container")

	return commander.NewCommand("link", "Link a host directory into a container",
		fs,
		func(args []string) error {
			if "" == *n {
				return errors.New("You must name the container into which you want to link a directory (-n param)")
			}
			if "" == *src {
				return errors.New("You must specify the source directory to link into the container (-src)")
			}
			if "" == *dest {
				return errors.New("You must specify the destination directory inside the container (-dest)")
			}
			return LinkDirIntoContainer(*n, *src, *dest)
		})
}
