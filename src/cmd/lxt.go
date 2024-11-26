package main

import (
	`fmt`
	"log"
	`os`
	`os/user`
	`reflect`
	`path/filepath`

	"github.com/alecthomas/kong"
	"github.com/craigmj/commander"
	`github.com/docopt/docopt.go`
	`gopkg.in/yaml.v3`

	"lxt/lxt"

)

func isNil(a any) bool {
	return reflect.ValueOf(a).IsNil()
}

func firstOf(a ...any) any {
	for _, b := range a {
		if !isNil(b) {
			return b
		}
	}
	return nil
}

func doError(err error) {
	if nil==err {
		return
	}
	log.Fatal(err.Error())
}

var cli struct {
	AnsibleFacts lxt.AnsibleFactsKong `cmd`
}



func main() {
	cwd, err := os.Getwd()
	doError(err)
	whoami, err := user.Current()
	doError(err)
usage := `lxt. A tool for simplifying lxc management.

Usage:
	lxt ansible-facts --name <name>
	lxt autostart --name <name> [--off]
	lxt cp <source> <destination>
	lxt config --name <name>
	lxt cpath <path>...
	lxt exists --name <name>
	lxt expose-port --name <name> --host-port <host-port> --container-port <container-port>
	lxt fix-ip --name <name> --ip <ip>
	lxt get-ip --name <name>
	lxt host --name <name> --source <source> --dest <dest> --script <script>
	lxt ln --name <name> [--source <source>] [--dest <dest>] [<source> <dest>]
	lxt portforward --name <name> --host-port <host-port> --container-port <container-port>
	lxt shell --name <name> <script>...
	lxt sudoer --name <name> [--user <user>]
	lxt tee [--append] <dest>...

Options:
	-h, --help  Show this help.
	-n, --name=<name>  Name of the container (default: ` + filepath.Base(cwd) + `).
	-a, --append  Append to the named file(s) (default: false).
	--source=<source>  Source directory to link (default: ` + cwd + `).
	--dest=<dest>  Destination directory in container (default: ` + filepath.Join(`/opt`, filepath.Base(cwd)) + `)
	--dir=<dir>  Directory to link into container (default: ` + cwd + `)
	--off  Turn feature off (such as autostart) (default: false).
	--container-port=<container-port>  Port on the container.
	--host-port=<host-port>  Port in the host.
	--ip=<ip>  IP address to fix for a container.
	--script=<script>  Script to execute in container (default: install.sh).
	--user=<user>  Username (default: ` + whoami.Username + `).

Description
== a description of lxt to come here ==

`	
	opts, err := docopt.ParseDoc(usage)
	doError(err)
	if false {
		yaml.NewEncoder(os.Stdout).Encode(opts)
		fmt.Println(`---`)
	}
	switch {
	case opts[`ansible-facts`]:
		doError(lxt.AnsibleFacts(opts[`--name`].(string)))
	case opts[`autostart`]:
		off, err := opts.Bool(`--off`)
		doError(err)
		doError(lxt.SetContainerAutoStart(opts[`--name`].(string), !off))
	case opts[`cp`]:
		doError(lxt.CopyFile(opts[`<source>`].(string), opts[`<destination>`].(string)))
	case opts[`config`]:
		doError(lxt.PrintLxtConfigFile(opts[`--name`].(string)))
	case opts[`cpath`]:
		doError(lxt.CPath(opts[`<path>`].([]string)))
	case opts[`exists`]:
		doError(lxt.ExistsCommandExit(opts[`--name`].(string)))
	case opts[`expose-port`]:
	case opts[`fix-ip`]:
		doError(lxt.FixIP(opts[`--name`].(string), opts[`--ip`].(string)))
	case opts[`get-ip`]:
		doError(lxt.PrintIP(opts[`--name`].(string)))
	case opts[`host`]:
		doError(lxt.InstallIntoContainer(
			opts[`--name`].(string),
			opts[`--source`].(string),
			opts[`--dest`].(string),
			opts[`--script`].(string),
		))
	case opts[`ln`]:
		doError(lxt.LinkCommand(opts[`--name`].(string), 
			firstOf(opts[`<source>`], opts[`--source`]).(string),
			firstOf(opts[`<dest>`], opts[`--dest`]).(string),
			[]string{},
		))
	case opts[`portforward`]:
		hostPort, err := opts.Int(`--host-port`)
		doError(err)
		containerPort, err := opts.Int(`--container-port`)
		doError(err)
		doError(lxt.PortForward(opts[`--name`].(string), hostPort, containerPort))
	case opts[`shell`]:
		for _, s := range opts[`<script>`].([]string) {
			doError(lxt.Shell(opts[`--name`].(string), s))
		}
	case opts[`sshkey`]:
		doError(lxt.PrintSSHKey(opts[`--name`].(string)))
	case opts[`ssh-keyscan`]:
		doError(lxt.SSHKeyScan(opts[`--name`].(string), opts[`user`].(string)))
	case opts[`sudoer`]:
		doError(lxt.SetSudoer(opts[`--name`].(string), opts[`--user`].(string)))
	case opts[`tee`]:
		append, err := opts.Bool(`--append`)
		doError(err)
		doError(lxt.Tee(append, opts[`<dest>`].([]string)))
	case opts[`write-ansible-fact`]:
		doError(lxt.WriteAnsibleFact(opts[`<name>`].(string)))
	}

	return 

	ctx := kong.Parse(&cli)
	switch ctx.Command() {
	default:
		doError(ctx.Run(ctx))
	}
	if err := commander.Execute(nil,
		lxt.AnsibleFactsCommand,
		lxt.AutoStartCommand,
		lxt.CopyCommand,
		lxt.CPathCommand,
		lxt.ExistsCommand,
		lxt.ExposePortCommand,
		lxt.FixIPCommand,
		lxt.GetIPCommand,
		lxt.HostCommand,
		lxt.LnCommand,
		lxt.PortForwardCommand,
		lxt.ShellCommand,
		lxt.SSHKeyCommand,
		lxt.SSHKeyScanCommand,
		lxt.SudoerCommand,
		lxt.TeeCommand,
		lxt.WriteAnsibleFactCommand,
	); nil != err {
		log.Fatal(err)
	}
}
