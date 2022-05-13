//go:build linux || freebsd || netbsd || openbsd
// +build linux freebsd netbsd openbsd

package file

import (
	"os"
	"syscall"
)

func (i info) ctime() int64 {
	stat := i.fd.Sys().(*syscall.Stat_t)
	return stat.Ctim.Nsec
}

func (i info) atime() int64 {
	stat := i.fd.Sys().(*syscall.Stat_t)
	return stat.Atim.Nsec
}

func openFile(filename string) (*os.File, error) {
	return os.Open(filename)
}
