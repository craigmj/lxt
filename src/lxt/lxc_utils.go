package lxt

import (
	"errors"

	"gopkg.in/lxc/go-lxc.v2"
)

func GetDefinedContainer(n string) (*lxc.Container, error) {
	if "" == n {
		return nil, errors.New("You need to provide a container name (-n)")
	}
	cont, err := lxc.NewContainer(n, lxc.DefaultConfigPath())
	if nil != err {
		return nil, err
	}
	if !cont.Defined() {
		return nil, errors.New("No container named " + n + " is defined. Did you forget to sudo?")
	}
	return cont, nil
}
