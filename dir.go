package file

import (
	"github.com/vela-security/vela-public/lua"
	"io/ioutil"
	"os"
	"path/filepath"
)

type dir struct {
	path string
	data []os.FileInfo
	err  error
}

func newLuaFileDir(L *lua.LState) int {
	path := L.CheckString(1)
	d, err := ioutil.ReadDir(path)
	L.Push(L.NewAnyData(&dir{path: path, data: d, err: err}))
	return 1
}

func (d *dir) ok() bool {
	if d.err != nil {
		return false
	}

	return true
}

//d.cb(=>(ev)

//end)

func (d *dir) visit(L *lua.LState, fn *lua.LFunction, m func(string) bool) int {
	if fn == nil {
		return 0
	}

	n := len(d.data)
	if n == 0 {
		return 0
	}

	co := xEnv.Clone(L)
	defer xEnv.Free(co)

	cp := xEnv.P(fn)
	for i := 0; i < n; i++ {
		v := d.data[i]
		uv := newInfo(filepath.Join(d.path, v.Name()), v, nil)
		if m != nil && !m(uv.path) {
			continue
		}
		err := co.CallByParam(cp, uv)
		if err != nil {
			xEnv.Errorf("call dir ipairs error %v", err)
			continue
		}
		co.SetTop(0)
	}
	return 0
}

func (d *dir) grep(L *lua.LState) int {
	pn := L.IsString(1)
	fn := L.IsFunc(2)

	if pn == "" {
		return d.visit(L, fn, nil)
	}

	return d.visit(L, fn, func(path string) bool {
		m, err := filepath.Match(pn, path)
		if err != nil {
			xEnv.Errorf("%s grep match %v", L.CodeVM(), err)
			return false
		}
		return m
	})

}

func (d *dir) ipairs(L *lua.LState) int {
	return d.visit(L, L.IsFunc(1), nil)
}

func (d *dir) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "ok":
		return lua.LBool(d.ok())
	case "err":
		if d.ok() {
			return lua.LNil
		}
		return lua.S2L(d.err.Error())

	case "count":
		return lua.LInt(len(d.data))
	case "grep":
		return L.NewFunction(d.grep)

	case "ipairs":
		return L.NewFunction(d.ipairs)
	}

	return lua.LNil
}
