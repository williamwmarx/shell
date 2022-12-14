#!/bin/sh
OS=`uname -s`
ARCH=`uname -m`

if [ ${OS} = Darwin ] || [ ${OS} = Linux ]; then
	curl -s https://releases.marx.sh/shell-config-${OS}-${ARCH} -o $HOME/shell-config
	chmod +x $HOME/shell-config
	$HOME/shell-config $@
	rm $HOME/shell-config
else
	echo "Shell config only works on macOS (arm64/amd64) and Linux (arm/arm64/amd64)"
fi
