package file

import "github.com/vela-security/vela-public/lua"

func Open(L *lua.LState, name string, path string, warp string) lua.Writer {
	cfg := &config{name: name, path: path, delim: warp}
	if e := cfg.verify(); e != nil {
		L.RaiseError("%v", e)
		return nil
	}

	xf := newFile(cfg)
	xf.V(lua.PTPrivate)
	if e := xf.Start(); e != nil {
		L.RaiseError("%v", e)
		return nil
	}

	return xf
}
