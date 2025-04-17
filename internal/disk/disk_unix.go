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

// GetInfo returns total and free bytes available in a directory, e.g. `/`.
func GetInfo(path string, firstTime bool) (info Info, err error) {
    s := syscall.Statfs_t{}
    err = syscall.Statfs(path, &s)
    if err != nil {
        return Info{}, err
    }
    reservedBlocks := s.Bfree - s.Bavail
    info = Info{
        Total: uint64(s.Frsize) * (s.Blocks - reservedBlocks),
        Free:  uint64(s.Frsize) * s.Bavail,
        Files: s.Files,
        Ffree: s.Ffree,
        //nolint:unconvert
        FSType: getFSType(uint32(s.Type)),
    }

    st := syscall.Stat_t{}
    err = syscall.Stat(path, &st)
    if err != nil {
        return Info{}, err
    }
    //nolint:unconvert
    devID := uint64(st.Dev) // Needed to support multiple GOARCHs
    info.Major = unix.Major(devID)
    info.Minor = unix.Minor(devID)

    // Check for overflows.
    // https://github.com/minio/minio/issues/8035
    // XFS can show wrong values at times error out
    // in such scenarios.
    if info.Free > info.Total {
        return info, fmt.Errorf("detected free space (%d) > total drive space (%d), fs corruption at (%s). please run 'fsck'", info.Free, info.Total, path)
    }
    info.Used = info.Total - info.Free

    if firstTime {
        bfs, err := blockdevice.NewDefaultFS()
        if err == nil {
            devName := ""
            diskstats, _ := bfs.ProcDiskstats()
            for _, dstat := range diskstats {
                // ignore all loop devices
                if strings.HasPrefix(dstat.DeviceName, "loop") {
                    continue
                }
                if dstat.MajorNumber == info.Major && dstat.MinorNumber == info.Minor {
                    devName = dstat.DeviceName
                    break
                }
            }
            if devName != "" {
                info.Name = devName
                qst, err := bfs.SysBlockDeviceQueueStats(devName)
                if err != nil { // Mostly not found error
                    // Check if there is a parent device:
                    //   e.g. if the mount is based on /dev/nvme0n1p1, let's calculate the
                    //        real device name (nvme0n1) to get its sysfs information
                    parentDevPath, e := os.Readlink("/sys/class/block/" + devName)
                    if e == nil {
                        parentDev := filepath.Base(filepath.Dir(parentDevPath))
                        qst, err = bfs.SysBlockDeviceQueueStats(parentDev)
                    }
                }
                if err == nil {
                    info.NRRequests = qst.NRRequests
                    rot := qst.Rotational == 1 // Rotational is '1' if the device is HDD
                    info.Rotational = &rot
                }
            }
        }
    }

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
