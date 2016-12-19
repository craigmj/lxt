package lxt

// import (
// 	"flag"
// 	"os"

// 	"github.com/craigmj/commander"
// )

// func SSHAuthKey(n string, u string, keyFile string) error {
// 	raw, err := ioutil.ReadFile(keyFile)
// 	if nil != err {
// 		return err
// 	}
// 	auth_keys := fmt.Sprintf(`%s:/home/%s/.ssh/authorized_keys`, n, u)
// 	auth_key_file, err := ContainerizePath(auth_keys)
// 	if nil != err {
// 		return err
// 	}
// 	out, err := os.OpenFile(auth_key_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
// 	if nil != err {
// 		return err
// 	}
// 	defer out.Close()
// 	defer

// }

// func SSHAuthKeyCommand() *commander.Command {
// 	fs := flag.NewFlagSet(`ssh-authkey`, flag.ExitOnError)
// 	n := fs.String(`n`, ``, `Name of the container`)
// 	u := fs.String(`u`, ``, `User in the container to whom the current user's key should be added`)
// 	key := fs.String(`key`, ``, `File containing the key to be authorized`)
// 	return commander.NewCommand(`ssh-authkey`,
// 		`Add an authorized SSH key for the given user in the container`,
// 		fs,
// 		func([]string) error {
// 			return SSHAuthKey(*n, *u, *key)
// 		})
// }
