package file

import (
	"fmt"
	"github.com/vela-security/vela-public/lua"
	"os"
	"path/filepath"
)

type info struct {
	path string
	fd   os.FileInfo
	ext  string
	err  error
}

func (i info) String() string                         { return fmt.Sprintf("%+v", &i) }
func (i info) Type() lua.LValueType                   { return lua.LTObject }
func (i info) AssertFloat64() (float64, bool)         { return 0, false }
func (i info) AssertString() (string, bool)           { return "", false }
func (i info) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (i info) Peek() lua.LValue                       { return i }

func newInfo(path string, fd os.FileInfo, err error) info {
	return info{
		path: path,
		fd:   fd,
		ext:  filepath.Ext(path),
		err:  err,
	}
}

func (i info) ok() bool {
	if i.err == nil {
		return true
	}
	return false
}

func (i info) Index(L *lua.LState, key string) lua.LValue {
	switch key {

	case "ok":
		return lua.LBool(i.ok())

	case "name":
		if i.ok() {
			return lua.LNil
		}

		return lua.S2L(i.fd.Name())

	case "ext":
		if i.ok() {
			return lua.LNil
		}
		return lua.S2L(i.ext)

	case "mtime":
		if i.ok() {
			return lua.LNil
		}
		return lua.LNumber(i.fd.ModTime().Unix())

	case "ctime":
		if i.ok() {
			return lua.LNil
		}
		return lua.LNumber(i.ctime())

	case "atime":
		if i.ok() {
			return lua.LNil
		}
		return lua.LNumber(i.atime())

	case "path":
		if i.ok() {
			return lua.LNil
		}
		return lua.S2L(i.path)

	case "dir":
		if i.ok() {
			return lua.LNil
		}
		return lua.LBool(i.fd.IsDir())

	}

	return lua.LNil
}
