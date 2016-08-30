package lxt

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/craigmj/commander"

	"gopkg.in/lxc/go-lxc.v2"
)

var _ = fmt.Printf
var _ = filepath.Join
var _ = os.DevNull

func bashScript(c *lxc.Container, workingDirectory string, script string) error {
	sin, sout, err := os.Pipe()
	if nil != err {
		return err
	}
	go func() {
		fmt.Fprintln(sout, `su - ubuntu`)
		fmt.Fprintln(sout, script) // /opt/bookworks/install.sh`)
		sout.Close()
	}()
	done, err := c.RunCommand([]string{
		`/bin/bash`,
	}, lxc.AttachOptions{
		UID:        -1,
		GID:        -1,
		Cwd:        workingDirectory,
		StdinFd:    sin.Fd(),
		StdoutFd:   os.Stdout.Fd(),
		StderrFd:   os.Stderr.Fd(),
		Namespaces: -1,
	})
	if nil != err {
		return err
	}
	if !done {
		return fmt.Errorf("Script %s exited with error", script)
	}
	return nil
}

func installIntoContainer(n string, hostDir, containerDir string, script string) error {
	c, err := lxc.NewContainer(n, lxc.DefaultConfigPath())
	if nil != err {
		return err
	}

	if true {

		if c.Defined() {
			return fmt.Errorf("Container %s already exists", n)
		}
		if err = c.Create(lxc.TemplateOptions{
			Template: `ubuntu`,
		}); nil != err {
			return err
		}

		// Link the current directory as /opt/bookworks
		out, err := os.OpenFile(c.ConfigFileName(), os.O_WRONLY|os.O_APPEND, 0664)
		if nil != err {
			return err
		}
		defer out.Close()
		fmt.Fprintf(out,
			"lxc.mount.entry = %s %s none bind,optional,create=dir 0 0\n",
			hostDir, stripLeadingSlash(containerDir))
		out.Close()

		// Setup ubuntu as a sudoer
		// fmt.Println("c.ConfigPath() = ", c.ConfigPath())
		sudoers, err := os.Create(filepath.Join(c.ConfigPath(), n, "rootfs",
			"etc", "sudoers.d", "ubuntu"))
		if nil != err {
			return err
		}
		defer sudoers.Close()
		fmt.Fprintln(sudoers, `ubuntu ALL=(ALL:ALL) NOPASSWD:ALL`)
		sudoers.Close()

		// Force a reload of the container config
		c, err = lxc.NewContainer(n, lxc.DefaultConfigPath())
		if nil != err {
			return err
		}
		if err = c.LoadConfigFile(c.ConfigFileName()); nil != err {
			return err
		}
		if err = c.Start(); nil != err {
			return err
		}
		fmt.Println("Waiting for container to run")
		if !c.Wait(lxc.RUNNING, 10*time.Second) {
			return fmt.Errorf("Failed to reach State in 10s")
		}
		if err := bashScript(c, containerDir, `
			while ! ping -c 1 -W 1 8.8.8.8; do
				echo 'Waiting for internet...'
				sleep 1
			done
		`); nil != err {
			return err
		}
	}

	if err := bashScript(c, containerDir, filepath.Join(containerDir, script)); nil != err {
		return err
	}
	return nil
}

func stripLeadingSlash(path string) string {
	if path[0] == '/' {
		return path[1:]
	}
	return path
}

func HostCommand() *commander.Command {
	fs := flag.NewFlagSet("host", flag.ExitOnError)
	cwd, err := os.Getwd()
	if nil != err {
		panic(err)
	}
	cname := fs.String("n", filepath.Base(cwd), "Name of the new container")
	dir := fs.String("dir", cwd, "Directory to link into container")
	dest := fs.String("dest", filepath.Join("/opt", filepath.Base(cwd)), "Destination directory in container")
	script := fs.String("script", "install.sh", "Installation script to execute in container")

	return commander.NewCommand(
		"host",
		"Install script into container",
		fs,
		func([]string) error {
			return installIntoContainer(*cname, *dir, *dest, *script)
		})
}
