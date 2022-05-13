package file

import (
	"github.com/vela-security/vela-public/auxlib"
	"github.com/vela-security/vela-public/lua"
	"os"
	"time"
)

func (xf *xFile) pushL(L *lua.LState) int {
	n := L.GetTop()
	if n == 0 {
		L.Push(lua.LNil)
		return 0
	}

	for i := 1; i <= n; i++ {
		xf.Push(auxlib.Format(L, 0))
	}
	return 0

}

func (xf *xFile) backupL(L *lua.LState) int {
	filename := xf.filename(time.Now())

	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		xEnv.Errorf("file backup fail %v", err)
		return 0
	}

	old := xf.fd
	defer old.Close()
	xf.fd = fd
	return 0
}

func (xf *xFile) Index(L *lua.LState, key string) lua.LValue {
	if key == "push" {
		return L.NewFunction(xf.pushL)
	}
	if key == "backup" {
		return L.NewFunction(xf.backupL)
	}

	return lua.LNil
}
