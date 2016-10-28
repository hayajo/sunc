// +build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

var defaultConfig = Config{
	Hostname: "sunc",
	Mounts: []Mount{
		{
			Source:     "tmpfs",
			Target:     "/tmp",
			Filesystem: "tmpfs",
			Flags:      syscall.MS_NOSUID | syscall.MS_NODEV,
		},
		{
			Source:     "tmpfs",
			Target:     "/run/lock",
			Filesystem: "tmpfs",
			Flags:      syscall.MS_NOSUID | syscall.MS_NODEV,
		},
		{
			Source:     "tmpfs",
			Target:     "/var/tmp",
			Filesystem: "tmpfs",
			Flags:      syscall.MS_NOSUID | syscall.MS_NODEV,
		},
		{
			Source:     "proc",
			Target:     "/proc",
			Filesystem: "proc",
			Flags:      syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_RELATIME,
		},
		{
			Source:     "tmpfs",
			Target:     "/dev",
			Filesystem: "tmpfs",
			Data:       "mode=755",
		},
		{
			Source:     "devpts",
			Target:     "/dev/pts",
			Filesystem: "devpts",
			Data:       "newinstance,ptmxmode=666",
		},
		{
			Source:     "/dev/pts/ptmx",
			Target:     "/dev/ptmx",
			Filesystem: "symlink",
		},
		{
			Source:     "/proc/self/fd",
			Target:     "/dev/fd",
			Filesystem: "symlink",
		},
		{
			Source:     "/proc/self/fd/1",
			Target:     "/dev/stdout",
			Filesystem: "symlink",
		},
		{
			Source:     "/proc/self/fd/2",
			Target:     "/dev/stderr",
			Filesystem: "symlink",
		},
		{
			Source:     "/proc/self/fd/0",
			Target:     "/dev/stdin",
			Filesystem: "symlink",
		},
		{
			Source:     "/dev/full",
			Target:     "/dev/full",
			Filesystem: "bind",
			Flags:      syscall.MS_BIND,
		},
		{
			Source:     "/dev/null",
			Target:     "/dev/null",
			Filesystem: "bind",
			Flags:      syscall.MS_BIND,
		},
		{
			Source:     "/dev/random",
			Target:     "/dev/random",
			Filesystem: "bind",
			Flags:      syscall.MS_BIND,
		},
		{
			Source:     "/dev/tty",
			Target:     "/dev/tty",
			Filesystem: "bind",
			Flags:      syscall.MS_BIND,
		},
		{
			Source:     "/dev/urandom",
			Target:     "/dev/urandom",
			Filesystem: "bind",
			Flags:      syscall.MS_BIND,
		},
		{
			Source:     "/dev/zero",
			Target:     "/dev/zero",
			Filesystem: "bind",
			Flags:      syscall.MS_BIND,
		},
		// {
		// Source:     "sysfs",
		// Target:     "/sys",
		// Filesystem: "sysfs",
		// Flags:      syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_RELATIME,
		// },
		// {
		// Source:     "/sys",
		// Target:     "/sys",
		// Filesystem: "rbind",
		// Flags:      syscall.MS_BIND | syscall.MS_REC,
		// },
		{
			Source:     "/etc/resolv.conf",
			Target:     "/etc/resolv.conf",
			Filesystem: "bind",
			Flags:      syscall.MS_BIND,
		},
		{
			Source:     "/etc/hosts",
			Target:     "/etc/hosts",
			Filesystem: "bind",
			Flags:      syscall.MS_BIND,
		},
	},
	Caps: []string{
		"CAP_AUDIT_WRITE",
		"CAP_CHOWN",
		"CAP_DAC_OVERRIDE",
		"CAP_FSETID",
		"CAP_FOWNER",
		"CAP_KILL",
		"CAP_MKNOD",
		"CAP_NET_BIND_SERVICE",
		"CAP_SETUID",
		"CAP_SETGID",
		"CAP_SETPCAP",
		"CAP_SETFCAP",
		"CAP_SYS_CHROOT",
	},
}

type Proc struct {
	Args   []string
	Config Config
}

func (p *Proc) Run() error {
	if err := p.init(); err != nil {
		return fmt.Errorf("Failed to init process: %v", err)
	}

	os.Clearenv()
	os.Setenv("CONTAINER", "sunc")

	name := p.Args[0]
	path, err := exec.LookPath(name)
	if err != nil {
		return err
	}

	return syscall.Exec(path, p.Args, os.Environ())
}

func (p *Proc) init() error {
	var err error

	err = p.setHostname()
	if err != nil {
		return err
	}

	err = p.mount()
	if err != nil {
		return err
	}

	err = p.chroot()
	if err != nil {
		return err
	}

	err = p.keepCaps()
	if err != nil {
		return err
	}

	return nil
}

func (p *Proc) mount() error {
	for _, m := range p.Config.Mounts {
		if err := m.Mount(p.Config.Rootfs); err != nil {
			return err
		}
	}
	return nil
}

func (p *Proc) setHostname() error {
	if p.Config.Hostname != "" {
		return syscall.Sethostname([]byte(p.Config.Hostname))
	}
	return nil
}

func (p *Proc) chroot() error {
	if p.Config.NoPivotRoot {
		return Chroot(p.Config.Rootfs)
	}
	return PivotRoot(p.Config.Rootfs)
}

func (p *Proc) keepCaps() error {
	return KeepCaps(os.Getpid(), p.Config.Caps)
}

func Process(name string, args ...string) *Proc {
	config := defaultConfig
	wd, _ := os.Getwd()
	config.Rootfs = wd
	return &Proc{
		Args:   append([]string{name}, args...),
		Config: config,
	}
}
