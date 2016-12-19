package lxt

import (
	"flag"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/craigmj/commander"
)

// SSHKeyScan adds the host keys for the named container to the
// current user's known hosts.
func SSHKeyScan(name string, euser string) error {
	userStruct, err := user.Lookup(euser)
	if nil != err {
		return err
	}
	uid, err := strconv.ParseInt(userStruct.Uid, 10, 64)
	if nil != err {
		return err
	}
	gid, err := strconv.ParseInt(userStruct.Gid, 10, 64)
	if nil != err {
		return err
	}
	ip, err := GetIP(name)
	if nil != err {
		return err
	}
	knownHostsFile := filepath.Join(`/home`, euser, `.ssh`, `known_hosts`)
	// If knownHostsFile exists, we remove all keys for the container
	_, err = os.Stat(knownHostsFile)
	if nil != err {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if err := exec.Command(`ssh-keygen`, `-R`, ip, `-f`, knownHostsFile).Run(); nil != err {
			return err
		}
	}
	scan := exec.Command(`ssh-keyscan`, `-H`, ip)
	app, err := os.OpenFile(knownHostsFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0600)
	if nil != err {
		return err
	}
	defer app.Close()
	defer os.Chown(knownHostsFile, int(uid), int(gid))
	scan.Stdout = app
	return scan.Run()
}

func SSHKeyScanCommand() *commander.Command {
	fs := flag.NewFlagSet(`ssh-keyscan`, flag.ExitOnError)
	n := fs.String("n", "", "Name of container")
	u := fs.String("u", "", "User")
	return commander.NewCommand(`ssh-keyscan`,
		`Add the host-keys for the named container to the host`,
		fs,
		func([]string) error {
			return SSHKeyScan(*n, *u)
		})
}
