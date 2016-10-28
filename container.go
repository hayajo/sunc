// +build linux

package main

import (
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

const INIT_COMMAND = "sunc-init"
const DEFAULT_FLAGS uintptr = syscall.CLONE_NEWUSER | syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS

type Ct struct {
	Args  []string
	Flags uintptr
	UID   int
	GID   int
}

func Container(name string, arg ...string) *Ct {
	ct := &Ct{
		Args:  append([]string{name}, arg...),
		UID:   os.Getuid(),
		GID:   os.Getgid(),
		Flags: DEFAULT_FLAGS,
	}
	return ct
}

func (ct *Ct) Start() error {
	cmd := exec.Command("/proc/self/exe")

	cmd.Args = append([]string{INIT_COMMAND}, ct.Args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	attr := &syscall.SysProcAttr{
		Cloneflags: ct.Flags,
	}

	if attr.Cloneflags&syscall.CLONE_NEWUSER > 0 {
		u, err := user.Current()
		if err != nil {
			return err
		}

		uid, _ := strconv.Atoi(u.Uid)
		uidMaps := []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      uid,
				Size:        1,
			},
		}

		gid, _ := strconv.Atoi(u.Gid)
		gidMaps := []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      gid,
				Size:        1,
			},
		}

		if os.Geteuid() == 0 { // suid root
			attr.GidMappingsEnableSetgroups = true
			attr.Credential = &syscall.Credential{
				Uid: 0,
				Gid: 0,
			}

			subuids, _ := GetSubUID(u)
			uidMaps = append(uidMaps, SubIDs(subuids).ToSysProcIDMaps(1)...)
			subgids, _ := GetSubGID(u)
			gidMaps = append(gidMaps, SubIDs(subgids).ToSysProcIDMaps(1)...)
		}

		attr.UidMappings = uidMaps
		attr.GidMappings = gidMaps
	}

	cmd.SysProcAttr = attr

	return cmd.Run()
}
