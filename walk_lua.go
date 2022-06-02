package file

import (
	"github.com/vela-security/vela-public/grep"
	"github.com/vela-security/vela-public/lua"
	"github.com/vela-security/vela-public/pipe"
)

//func (w *walk) String() string                         { return fmt.Sprintf("%p", w) }
//func (w *walk) Type() lua.LValueType                   { return lua.LTObject }
//func (w *walk) AssertFloat64() (float64, bool)         { return 0, false }
//func (w *walk) AssertString() (string, bool)           { return "", false }
//func (w *walk) AssertFunction() (*lua.LFunction, bool) { return nil, false }
//func (w *walk) Peek() lua.LValue                       { return w }

//func (w *walk) setExt(L *lua.LState) int {
//	n := L.GetTop()
//	if n == 0 {
//		return 0
//	}
//
//	w.ext = make([]string, n)
//	for i := 1; i <= n; i++ {
//		w.ext[i-1] = L.CheckString(i)
//	}
//	return 0
//}
//
//func (w *walk) setNotExt(L *lua.LState) int {
//	n := L.GetTop()
//	if n == 0 {
//		return 0
//	}
//
//	w.notExt = make([]string, n)
//	for i := 1; i <= n; i++ {
//		w.ext[i-1] = L.CheckString(i)
//	}
//	return 0
//}

func (w *walk) dirL(L *lua.LState) int {
	w.cfg.dir = L.IsTrue(1)
	return 0
}

func (w *walk) limitL(L *lua.LState) int {
	w.cfg.limit = L.IsInt(1)
	return 0
}

func (w *walk) run(L *lua.LState) int {
	xEnv.Start(L, w).From(L.CodeVM()).Do()
	return 0
}

func (w *walk) ignoreL(L *lua.LState) int {
	n := L.GetTop()
	if n == 0 {
		return 0
	}
	for i := 1; i <= n; i++ {
		w.cfg.ignore = append(w.cfg.ignore, grep.New(L.IsString(i)))
	}
	return 0
}

func (w *walk) filterL(L *lua.LState) int {
	n := L.GetTop()
	if n == 0 {
		return 0
	}
	for i := 1; i <= n; i++ {
		w.cfg.filter = append(w.cfg.filter, grep.New(L.IsString(i)))
	}
	return 0
}

func (w *walk) pipeL(L *lua.LState) int {
	w.cfg.pipe.CheckMany(L, pipe.Seek(0))
	return 0
}

func (w *walk) deepL(L *lua.LState) int {
	n := L.IsInt(1)
	if n == 0 {
		return 0
	}

	w.cfg.deep = n
	return 0
}

func (w *walk) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "dir":
		return L.NewFunction(w.dirL)
	case "deep":
		return L.NewFunction(w.deepL)

	case "ignore":
		return lua.NewFunction(w.ignoreL)
	case "filter":
		return lua.NewFunction(w.filterL)
	case "pipe":
		return lua.NewFunction(w.pipeL)
	//case "ext":
	//	return L.NewFunction(w.setExt)
	//case "not_ext":
	//	return L.NewFunction(w.setNotExt)
	//case "wait":
	//	return L.NewFunction(w.wait)

	case "limit":
		return L.NewFunction(w.limitL)
	case "run":
		return L.NewFunction(w.run)
	case "zip":
		return L.NewFunction(w.zip)

	}
	return lua.LNil
}
