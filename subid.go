// +build linux

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

type IDType int

const (
	UID IDType = iota
	GID
)

var (
	subUIDFile = "/etc/subuid"
	subGIDFile = "/etc/subgid"
)

type SubID struct {
	ID    int
	Count int
}

func GetSubUID(u *user.User) (SubIDs, error) {
	return subIDs(UID, u)
}

func GetSubGID(u *user.User) (SubIDs, error) {
	return subIDs(GID, u)
}

func subIDs(t IDType, u *user.User) (SubIDs, error) {
	var filename string
	var subIDs SubIDs

	switch t {
	case UID:
		filename = subUIDFile
	case GID:
		filename = subGIDFile
	default:
		return subIDs, fmt.Errorf("Unknown ID type: %v", t)
	}

	file, err := os.Open(filename)
	if err != nil {
		return subIDs, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry := strings.SplitN(scanner.Text(), ":", 3)
		if len(entry) != 3 || entry[0] != u.Username {
			continue
		}

		id, err := strconv.Atoi(entry[1])
		if err != nil {
			continue
		}

		count, err := strconv.Atoi(entry[2])
		if err != nil {
			continue
		}

		subIDs = append(subIDs, SubID{
			ID:    id,
			Count: count,
		})
	}

	if err := scanner.Err(); err != nil {
		return subIDs, err
	}

	return subIDs, nil
}

type SubIDs []SubID

func (s SubIDs) Len() int           { return len(s) }
func (s SubIDs) Less(i, j int) bool { return s[i].ID < s[j].ID }
func (s SubIDs) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SubIDs) Sort()              { sort.Sort(s) }

func (s SubIDs) ToSysProcIDMaps(start int) []syscall.SysProcIDMap {
	var idMaps []syscall.SysProcIDMap

	cID := start
	s.Sort()
	for _, v := range s {
		idMaps = append(idMaps, syscall.SysProcIDMap{
			ContainerID: cID,
			HostID:      v.ID,
			Size:        v.Count,
		})

		cID += v.Count
	}

	return idMaps
}
