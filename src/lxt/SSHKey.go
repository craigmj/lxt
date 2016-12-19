package lxt

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/craigmj/commander"
)

// GetSSHKey gets the public SSH key for the given container
func GetSSHKey(name string) (string, error) {
	cont, err := GetDefinedContainer(name)
	if nil != err {
		return "", err
	}

	sshPath := filepath.Join(cont.ConfigPath(), name, "rootfs", "etc", "ssh", "ssh_host_ecdsa_key.pub")
	raw, err := ioutil.ReadFile(sshPath)
	if nil != err {
		return "", err
	}
	return string(raw), nil
}

func SSHKeyCommand() *commander.Command {
	fs := flag.NewFlagSet("sshkey", flag.ExitOnError)
	n := fs.String("n", "", "Name of the container")
	return commander.NewCommand("sshkey",
		"Get the public sshkey for the container",
		fs,
		func([]string) error {
			sshKey, err := GetSSHKey(*n)
			if nil != err {
				return err
			}
			fmt.Print(sshKey)
			return nil
		})
}
