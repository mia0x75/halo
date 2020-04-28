#!/bin/bash
git="$(git log --date=iso --pretty=format:"%cd" -1) $(git describe --tags --always)"
version=$(cat VERSION)
build=$(cat BUILD)
echo $(($(cat BUILD) + 1)) > BUILD
kernel=$(uname -r)
distro="Unknown"
os=$(uname | tr '[:upper:]' '[:lower:]')
case ${os} in
	linux*)
		name=$(cat /etc/*-release | tr [:upper:] [:lower:] | grep -Poi '(debian|ubuntu|red hat|centos|fedora)'|uniq)
		if [ ! -z $name ]; then
			distro=$(cat /etc/${name}-release)
		fi
		;;
	darwin*)
		distro="$(sw_vers -productName) $(sw_vers -productVersion) $(sw_vers -buildVersion)"
		;;
	*)
		;;
esac

if [ "X${git}" == "X" ]; then
    git="not a git repo"
fi

compile="$(date +"%F %T %z") by $(go version)"

branch=$(git rev-parse --abbrev-ref HEAD)

cat <<EOF | gofmt >g/g.go
package g

import (
	"runtime"
)

const (
	Version = "${version}.${build}"
	Git     = "${git}"
	Compile = "${compile}"
	Branch  = "${branch}"
	Distro  = "${distro}"
	Kernel  = "${kernel}"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
EOF
