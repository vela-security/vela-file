package file

import (
	"github.com/vela-security/vela-public/lua"
	"github.com/vela-security/vela-public/pipe"
	"io/fs"
	"io/ioutil"
	"path/filepath"
)

type dir struct {
	path   string
	data   []fs.FileInfo
	filter []func(string) bool
	pipe   *pipe.Px
	err    error
}

func newLuaFileDir(L *lua.LState) int {
	path := L.CheckString(1)
	data, err := ioutil.ReadDir(path)
	L.Push(&dir{
		path: path,
		data: data,
		err:  err,
		pipe: pipe.New(pipe.Env(xEnv)),
	})
	return 1
}

func (d *dir) ok() bool {
	if d.err != nil {
		return false
	}

	return true
}

func (d *dir) fuzzy() func(string) bool {
	if len(d.filter) == 0 {
		return func(_ string) bool {
			return true
		}
	}

	return func(path string) bool {
		return match(path, d.filter)
	}
}

func (d *dir) Info(idx int) info {
	fd := d.data[idx]
	path := filepath.Join(d.path, fd.Name())
	return newInfo(path, fd, nil)
}
