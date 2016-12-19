package lxt

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/craigmj/commander"
)

func WriteAnsibleFact() *commander.Command {
	fs := flag.NewFlagSet("write-ansible-fact", flag.ExitOnError)
	n := fs.String("n", "", "Name of container")
	return commander.NewCommand(
		`write-ansible-fact`,
		`writes an ansible fact to /etc/ansible/facts.d/[[container-name]].fact`,
		fs,
		func(args []string) error {
			ip, err := GetIP(*n)
			if nil != err {
				return err
			}
			dest := fmt.Sprintf("/etc/ansible/facts.d")
			os.MkdirAll(dest, 0755)
			out, err := os.Create(filepath.Join(dest, fmt.Sprintf("%s.fact", *n)))
			if nil != err {
				return err
			}
			defer out.Close()
			if err = json.NewEncoder(out).Encode(map[string]string{
				"ip": ip,
			}); nil != err {
				return err
			}

			return nil
		})
}

func AnsibleFacts() *commander.Command {
	fs := flag.NewFlagSet("ansible-facts", flag.ExitOnError)
	n := fs.String("n", "", "Name of container")
	return commander.NewCommand(
		`ansible-facts`,
		`Outputs some facts of the named container in ansible_facts JSON`,
		fs,
		func(args []string) error {
			ip, err := GetIP(*n)
			if nil != err {
				return err
			}

			sshKey, err := GetSSHKey(*n)
			if nil != err {
				return err
			}

			facts := map[string]interface{}{
				"ip":      ip,
				"ssh_key": sshKey,
			}

			if err = json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
				"ansible_facts": map[string]interface{}{
					"lxc": map[string]interface{}{
						*n: facts,
					},
				},
			}); nil != err {
				return err
			}

			return nil
		})
}
