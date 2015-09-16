# lxt
Simple tool for doing basic linux container tasks.

This exists to facilitate working with containers via scripts.

## link

Link links a source directory in the host to a destination directory in the container.

    sudo lxt link -n container_name -src /home/craigmj -dest opt/home

The above command links the /home/craigmj directory in the host to the opt/home directory in the container.

## exists

Exists checks whether a container with the given name exists. It returns to the shell a 0 if the container exists, or a 1 if it does not. So it can be used in shell scripts

	if sudo lxt exists -n test; then
		echo The container named test exists
	fi

