package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

func Chroot(root string) error {
	if err := syscall.Chroot(root); err != nil {
		return err
		return fmt.Errorf("Failed to chroot %s: %v", root, err)
	}
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("Failed to chdir /: %v", err)
	}
	return nil
}

const OLD_ROOT = ".old_root"

func PivotRoot(root string) error {
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Failed to mount root to itself: %v", root, err)
	}

	oldRoot := filepath.Join(root, OLD_ROOT)
	if err := os.Mkdir(oldRoot, 0777); err != nil {
		return fmt.Errorf("Failed to mkdir %s: %v", oldRoot, err)
	}
	if err := syscall.PivotRoot(root, oldRoot); err != nil {
		return fmt.Errorf("Failed to pivot_root: %v", err)
	}
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("Failed to chdir /: %v", err)
	}

	oldRoot = filepath.Join("/", OLD_ROOT)
	if err := syscall.Unmount(oldRoot, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("Failed to unmount old_root: %v", err)
	}
	if err := os.Remove(oldRoot); err != nil {
		return fmt.Errorf("Failed to remove old_root: %v", err)
	}

	return nil
}
