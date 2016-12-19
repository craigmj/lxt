package lxt

import (
	"errors"
	"flag"
	"net"
	"strconv"

	"github.com/craigmj/commander"
)

// ExposePort is going to expose a containerPort on the Host.
// The basic iptables commands to accomplish this are:
// #sudo iptables -N PHPC
// #sudo iptables -A PHPC -p tcp -d 10.0.3.192 --dport 80 -j ACCEPT
// #sudo iptables -A PHPC -j DROP
// #sudo iptables -A FORWARD -i eth0 -o lxcbr0 -j PHPC
// #sudo iptables -A FORWARD -j DROP	# not sure if this is necessary
// sudo iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 80 -j DNAT --to 10.0.3.192:80
//
func ExposePort(n string, containerPort, hostPort int) error {
	iface, err := DefaultInterface()
	if nil != err {
		return err
	}
	containerIp, err := GetIP(n)
	if nil != err {
		return err
	}
	action := "-C" // CHECK first

	ruleExists, err := IPTables(
		`-t`, `nat`,
		action, `PREROUTING`,
		`-i`, iface,
		`-p`, `tcp`,
		`--dport`, strconv.FormatInt(int64(hostPort), 10),
		`-j`, `DNAT`,
		`--to`, net.JoinHostPort(containerIp, strconv.FormatInt(int64(hostPort), 10)))
	if nil != err {
		// It apears that the CHECK can sometimes fail if nothing
		// has yet been added to PREROUTING, so we'll take an ERROR
		// as a possible 'rule doesn't yet exist'.
		ruleExists = false
	}
	if !ruleExists {
		// Rule does not exist
		action = "-A" // ADD the rule
		_, err := IPTables(
			`-t`, `nat`,
			action, `PREROUTING`,
			`-p`, `tcp`,
			`-i`, iface,
			`--dport`, strconv.FormatInt(int64(hostPort), 10),
			`-j`, `DNAT`,
			`--to`, net.JoinHostPort(containerIp, strconv.FormatInt(int64(hostPort), 10)))
		if nil != err {
			return err
		}
	}

	// Now we need to make the rule persistent
	return SaveIPTablesRules()
}

func ExposePortCommand() *commander.Command {
	fs := flag.NewFlagSet("expose-port", flag.ExitOnError)
	cport := fs.Int("container-port", 0, "Port insider container")
	hport := fs.Int("host-port", 0, "Port on host")
	n := fs.String(`n`, ``, `Container name`)
	return commander.NewCommand(`expose-port`,
		`Expose a container port in the named container as the specified port on the host`,
		fs,
		func([]string) error {
			if 0 == *cport {
				return errors.New("You need to specify the container port (-container-port)")
			}
			if 0 == *hport {
				return errors.New("You need to specify the host port (-host-port)")
			}
			return ExposePort(*n, *cport, *hport)
		})
}
