//go:build darwin || dragonfly
// +build darwin dragonfly

package disk

import (
	"fmt"
	"syscall"
)

// GetInfo returns total and free bytes available in a directory, e.g. `/`.
func GetInfo(path string, _ bool) (info Info, err error) {
	s := syscall.Statfs_t{}
	err = syscall.Statfs(path, &s)
	if err != nil {
		return Info{}, err
	}
	reservedBlocks := s.Bfree - s.Bavail
	info = Info{
		Total:  uint64(s.Bsize) * (s.Blocks - reservedBlocks),
		Free:   uint64(s.Bsize) * s.Bavail,
		Files:  s.Files,
		Ffree:  s.Ffree,
		FSType: getFSType(s.Fstypename[:]),
	}
	if info.Free > info.Total {
		return info, fmt.Errorf("detected free space (%d) > total drive space (%d), fs corruption at (%s). please run 'fsck'", info.Free, info.Total, path)
	}
	info.Used = info.Total - info.Free
	return info, nil
}

// getFSType returns the filesystem type of the underlying mounted filesystem
func getFSType(fstype []int8) string {
	b := make([]byte, len(fstype))
	for i, v := range fstype {
		b[i] = byte(v)
	}
	return string(b)
}
