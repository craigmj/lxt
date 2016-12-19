package lxt

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
)

var ErrNoDefaultInterfaceFound = errors.New(`Failed to find default interface`)

// DefaultInterface returns the name of the default interface
// on this system (eg 'eth0')
func DefaultInterface() (string, error) {
	route, err := exec.LookPath(`route`)
	if nil != err {
		return ``, fmt.Errorf(`Failed to find route command: %s`, err.Error())
	}
	cmd := exec.Command(route, `-n`)
	cmdOut, err := cmd.Output()
	if nil != err {
		return ``, fmt.Errorf(`Error running route -n: %s`, err.Error())
	}
	scan := bufio.NewScanner(bytes.NewReader(cmdOut))
	reg := regexp.MustCompile(
		`^0\.0\.0\.0\s+\S+\s+0\.0\.0\.0\s+\S+\s+\d+\s+\S+\s+\S+\s+(\S+)\s*$`)
	for scan.Scan() {
		// fmt.Printf("%s\n", scan.Text())
		if m := reg.FindStringSubmatch(scan.Text()); nil != m {
			return m[1], nil
		}
	}
	return ``, ErrNoDefaultInterfaceFound
}

// EnsureChain ensures that the chain exists, creating it
// if it does not.
func EnsureChain(chain string) error {
	iptables, err := exec.LookPath(`iptables`)
	if nil != err {
		return err
	}
	c := exec.Command(
		iptables, `-L`, chain)
	if err = c.Run(); nil != err {
		return err
	}
	if !c.ProcessState.Success() {
		return exec.Command(`iptables`,
			`-N`, chain).Run()
	}
	return nil
}

func SaveIPTablesRules() error {
	if err := exec.Command(
		`/etc/init.d/iptables-persistent`, `save`).Run(); nil != err {
		return fmt.Errorf("ERROR on /etc/init.d/iptables-persistent save: %s", err.Error())
	}
	return nil
}

func IPTables(args ...string) (success bool, err error) {
	iptables, err := exec.LookPath(`iptables`)
	if nil != err {
		return false, err
	}
	cmd := exec.Command(iptables, args...)
	if err := cmd.Run(); nil != err {
		return false, fmt.Errorf("ERROR on iptables %v: %s", args, err.Error())
	}
	return cmd.ProcessState.Success(), nil
}
