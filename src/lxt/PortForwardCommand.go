package lxt

import (
	"flag"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/craigmj/commander"
)

func PortForward(containerName string, portToForward int, destinationPort int) error {
	ip, err := GetIP(containerName)
	if nil != err {
		return err
	}
	/*
	 iptables
	 -A PREROUTING -
	 t nat
	 -p tcp
	 --dport 8000
	 -j DNAT
	 --to-destination YY.YY.YY.YY:8000
	*/

	iptables := exec.Command(
		`iptables`,
		`-A`, `PREROUTING`,
		`-t`, `nat`,
		`-p`, `tcp`,
		`--dport`, strconv.Itoa(portToForward),
		`-j`, `DNAT`,
		`--to-destination`, fmt.Sprintf("%s:%d", ip, destinationPort),
	)
	if err := iptables.Run(); nil != err {
		return err
	}
	iptables = exec.Command(
		`iptables`,
		`-A`, `FORWARD`,
		`-p`, `tcp`,
		`--sport`, strconv.Itoa(portToForward),
		`-d`, ip,
		`--dport`, strconv.Itoa(destinationPort),
		`-j`, `ACCEPT`,
	)
	if err := iptables.Run(); nil != err {
		return err
	}
	return nil
}

func PortForwardCommand() *commander.Command {
	fs := flag.NewFlagSet("portforward", flag.ExitOnError)
	n := fs.String("n", "", "Name of container")
	s := fs.Int("s", 8080, "Source Port to forward")
	d := fs.Int("d", 80, "Destination Port on container")

	return commander.NewCommand("portforward",
		"Forward a port on the host to the container",
		fs,
		func(args []string) error {
			return PortForward(*n, *s, *d)
		})
}
