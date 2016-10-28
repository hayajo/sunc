// +build linux

package main

type Config struct {
	Hostname    string
	Rootfs      string
	Mounts      []Mount
	Caps        []string
	NoPivotRoot bool
}
