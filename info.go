package file

import (
	"github.com/vela-security/vela-public/auxlib"
	"github.com/vela-security/vela-public/kind"
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

func (i info) String() string                         { return auxlib.B2S(i.Byte()) }
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

func (i info) Byte() []byte {
	enc := kind.NewJsonEncoder()
	enc.Tab("")
	enc.KV("path", i.path)
	enc.KV("ext", i.ext)
	enc.KV("mtime", i.MTime())
	enc.KV("size", i.fd.Size())
	enc.End("}")
	return enc.Bytes()
}

func (i info) ok() bool {
	return i.err == nil

}

func (i info) MTime() int64 {
	if i.ok() {
		return i.fd.ModTime().Unix()
	}
	return 0
}

func (i info) Index(L *lua.LState, key string) lua.LValue {
	switch key {

	case "name":
		return lua.S2L(i.fd.Name())

	case "ext":
		return lua.S2L(i.ext)

	case "mtime":
		return lua.LNumber(i.MTime())

	case "ctime":
		return lua.LNumber(i.ctime())

	case "atime":
		return lua.LNumber(i.atime())

	case "path":
		return lua.S2L(i.path)

	case "dir":
		return lua.LBool(i.fd.IsDir())

	}

	return lua.LNil
}
