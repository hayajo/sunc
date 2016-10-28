// +build linux

package main

import (
	"fmt"
	"strings"

	"github.com/syndtr/gocapability/capability"
)

var capMap map[string]capability.Cap

func init() {
	capMap = make(map[string]capability.Cap)

	last := capability.CAP_LAST_CAP

	// CentOS6対応
	// CentOS6では/proc/sys/kernel/cap_last_capが存在しないためCAP_BLOCK_SUSPENDはデフォルトの63のままとなる。
	// CentOS6での最後のcapは36(CAP_BLOCK_SUSPEND)となる。
	if last == capability.Cap(63) {
		last = capability.CAP_BLOCK_SUSPEND
	}

	for _, cap := range capability.List() {
		if cap > last {
			continue
		}
		capKey := fmt.Sprintf("CAP_%s", strings.ToUpper(cap.String()))
		capMap[capKey] = cap
	}
}

func KeepCaps(pid int, caps []string) error {
	var keep []capability.Cap
	for _, cap := range caps {
		v, ok := capMap[strings.ToUpper(string(cap))]
		if !ok {
			return fmt.Errorf("invalid capability %s", cap)
		}
		keep = append(keep, v)
	}

	cap, err := capability.NewPid(pid)
	if err != nil {
		return err
	}

	capType := capability.CAPS | capability.BOUNDING
	cap.Clear(capType)
	cap.Set(capType, keep...)

	return cap.Apply(capType)
}
