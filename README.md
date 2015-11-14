# lxt
Simple tool for doing basic linux container tasks.

This exists to facilitate working with containers via scripts.


## autostart

Autostart sets autostart on (or off) on the container.

	sudo lxt autostart -n CONTAINER

	or

	sudo lxt autostart -n CONTAINER -off

## cpath

Convert a Container path to a host path.

A container path is specified as CNAME:/path/in/container

	sudo lxt cpath test:/var/www

## Copy (cp)

Copy copies files between containers, or container and hosts. Use the CNAME:/path/in/container format to specify containers.

	sudo lxt cp /var/www/index.html test:/var/www/index.html

## exists

Exists checks whether a container with the given name exists. It returns to the shell a 0 if the container exists, or a 1 if it does not. So it can be used in shell scripts

	if sudo lxt exists -n test; then
		echo The container named test exists
	fi

## getip

Getip returns the IPv4 address of the container.

	sudo lxt getip -n CNAME

This can be useful for scripting, but be careful - it simply returns the FIRST ipv4 address listed for the container, even if there is more than one.

## Link (ln)

Link links a source directory in the host to a destination directory in the container.

    sudo lxt ln -n container_name -src /home/craigmj -dest opt/home

The above command links the /home/craigmj directory in the host to the opt/home directory in the container.

## portforward

Portfoward forwards all tcp ports on the host to the destination port on the container. It's really just a wrapper around a couple of iptables calls. Note, then, that if you want to maintain this configuration, you need to coordinate with saving your iptables configuration between reboots (see iptables-save or iptables-restore, for eg).

	sudo lxt portforward -s SOURCE_PORT -n CNAME -d DESTINATION_PORT

The IPtables configuration set by portforward uses the current ipv4 address of the container. Unless you are sure that this will be the same every time the container starts, you might _not_ want to save the iptables config...

## shell

Shell runs one or more shell scripts inside the container. It echos stdout and stderr to the hosts stdout and stderr. Shell is much easier than copying shell files into the container, then using lxc-execute...

	sudo lxt shell -n CNAME first.sh second.sh

## tee

Tee works just like Unix tee command (I think), copying stdin into one or more destination files. Use -a to append to those files instead of overwriting.

	cat test.txt | sudo lxt tee CNAME:/var/ww/test.txt


