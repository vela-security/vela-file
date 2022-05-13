//go:build darwin
// +build darwin

package file

import (
	"os"
	"syscall"
)

func (i info) ctime() int64 {
	stat := i.fd.Sys().(*syscall.Stat_t)
	return stat.Ctimespec.Nsec
}

func (i info) atime() int64 {
	stat := i.fd.Sys().(*syscall.Stat_t)
	return stat.Atimespec.Nsec
}

func openFile(filename string) (*os.File, error) {
	return os.Open(filename)
}
