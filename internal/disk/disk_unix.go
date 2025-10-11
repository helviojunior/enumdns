//go:build linux
// +build linux

package disk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"strconv"

	"github.com/prometheus/procfs/blockdevice"
	"golang.org/x/sys/unix"
)

// fsType2StringMap - list of filesystems supported on linux
var fsType2StringMap = map[string]string{
	"1021994":  "TMPFS",
	"137d":     "EXT",
	"4244":     "HFS",
	"4d44":     "MSDOS",
	"52654973": "REISERFS",
	"5346544e": "NTFS",
	"58465342": "XFS",
	"61756673": "AUFS",
	"6969":     "NFS",
	"ef51":     "EXT2OLD",
	"ef53":     "EXT4",
	"f15f":     "ecryptfs",
	"794c7630": "overlayfs",
	"2fc12fc1": "zfs",
	"ff534d42": "cifs",
	"53464846": "wslfs",
}

func getFilesystemStats(path string) (syscall.Statfs_t, error) {
	s := syscall.Statfs_t{}
	err := syscall.Statfs(path, &s)
	return s, err
}

func safeConvertStatfs(s syscall.Statfs_t) (frsize, blocks, bavail, ftype uint64) {
	reservedBlocks := s.Bfree - s.Bavail

	if s.Frsize >= 0 {
		frsize = uint64(s.Frsize)
	}
	if s.Blocks > reservedBlocks {
		blocks = uint64(s.Blocks - reservedBlocks)
	}
	if s.Bavail > 0 {
		bavail = uint64(s.Bavail)
	}
	if s.Type >= 0 {
		ftype = uint64(s.Type)
	}
	return
}

func getDeviceInfo(path string) (major, minor uint32, err error) {
	st := syscall.Stat_t{}
	err = syscall.Stat(path, &st)
	if err != nil {
		return 0, 0, err
	}
	devID := uint64(st.Dev)
	return unix.Major(devID), unix.Minor(devID), nil
}

func findDeviceName(major, minor uint32, bfs blockdevice.FS) string {
	diskstats, _ := bfs.ProcDiskstats()
	for _, dstat := range diskstats {
		if strings.HasPrefix(dstat.DeviceName, "loop") {
			continue
		}
		if dstat.MajorNumber == major && dstat.MinorNumber == minor {
			return dstat.DeviceName
		}
	}
	return ""
}

func processBlockDevice(info *Info, firstTime bool) error {
	if !firstTime {
		return nil
	}

	bfs, err := blockdevice.NewDefaultFS()
	if err != nil {
		return nil // Not a critical error
	}

	devName := findDeviceName(info.Major, info.Minor, bfs)
	if devName == "" {
		return nil
	}

	info.Name = devName
	qst, err := bfs.SysBlockDeviceQueueStats(devName)
	if err != nil {
		// Try parent device
		parentDevPath, e := os.Readlink("/sys/class/block/" + devName)
		if e == nil {
			parentDev := filepath.Base(filepath.Dir(parentDevPath))
			qst, err = bfs.SysBlockDeviceQueueStats(parentDev)
		}
	}

	if err == nil {
		info.NRRequests = qst.NRRequests
		rot := qst.Rotational == 1
		info.Rotational = &rot
	}

	return nil
}

// GetInfo returns total and free bytes available in a directory, e.g. `/`.
func GetInfo(path string, firstTime bool) (info Info, err error) {
	s, err := getFilesystemStats(path)
	if err != nil {
		return Info{}, err
	}

	frsize, blocks, bavail, ftype := safeConvertStatfs(s)

	info = Info{
		Total:  frsize * blocks,
		Free:   frsize * bavail,
		Files:  s.Files,
		Ffree:  s.Ffree,
		FSType: getFSType(uint32(ftype)),
	}

	info.Major, info.Minor, err = getDeviceInfo(path)
	if err != nil {
		return Info{}, err
	}

	if info.Free > info.Total {
		return info, fmt.Errorf("detected free space (%d) > total drive space (%d), fs corruption at (%s). please run 'fsck'", info.Free, info.Total, path)
	}
	info.Used = info.Total - info.Free

	processBlockDevice(&info, firstTime)

	return info, nil
}

// getFSType returns the filesystem type of the underlying mounted filesystem
func getFSType(ftype uint32) string {
	fsTypeHex := strconv.FormatUint(uint64(ftype), 16)
	fsTypeString, ok := fsType2StringMap[fsTypeHex]
	if !ok {
		return "UNKNOWN"
	}
	return fsTypeString
}
