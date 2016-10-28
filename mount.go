// +build linux

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

type Mount struct {
	Source     string
	Target     string
	Filesystem string
	Flags      uintptr
	Data       string
}

func (m *Mount) Mount(root string) error {
	target := filepath.Join(root, m.Target)
	var isDir bool

	switch m.Filesystem {
	case "symlink":
		return symlink(m.Source, target)
	case "bind":
		stat, err := os.Stat(m.Source)
		if err != nil {
			return err
		}
		isDir = stat.IsDir()
	default:
		isDir = true
	}

	if err := createTargetIfNotExists(target, isDir); err != nil {
		return fmt.Errorf("Failed to create mountpoint: %v", err)
	}

	return syscall.Mount(m.Source, target, m.Filesystem, m.Flags, m.Data)
}

func symlink(source, target string) error {
	if err := createTargetIfNotExists(filepath.Dir(target), true); err != nil {
		return fmt.Errorf("Failed to create mountpoint: %v", err)
	}
	return os.Symlink(source, target)
}

func createTargetIfNotExists(path string, isDir bool) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if isDir {
				err = os.MkdirAll(path, 0755)
				return err
			}
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(path, os.O_CREATE, 0755)
			if err != nil {
				return err
			}
			f.Close()
		}
	}
	return nil
}
