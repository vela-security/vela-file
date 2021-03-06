package file

import (
	"github.com/vela-security/vela-public/assert"
	"github.com/vela-security/vela-public/lua"
	"os"
)

var (
	xEnv assert.Environment
)

/*
	local w = file.open{name , path , delim}
	local st = file.stat(path)
	local wk = file.walk("name")

	local wx = wk.open("/var/log")
	wx.ext(".zip" , ".txt")
	wx.limit()
	wx.run()

*/

func newLuaFileOpen(L *lua.LState) int {
	cfg := newConfig(L)
	var ov *xFile

	proc := L.NewProc(cfg.name, fileTypeOf)
	if proc.IsNil() {
		ov = newFile(cfg)
		proc.Set(ov)
	} else {
		ov = proc.Data.(*xFile)
		ov.cfg = cfg
	}

	xEnv.Start(L, ov).From(L.CodeVM()).Do()
	L.Push(proc)
	return 1
}

func newLuaFileStat(L *lua.LState) int {
	path := L.IsString(1)
	if path == "" {
		return 0
	}

	fd, err := os.Stat(path)
	L.Push(newInfo(path, fd, err))
	return 1
}

func newLuaFileWalk(L *lua.LState) int {
	cfg := newWalkConfig(L)
	proc := L.NewProc(cfg.name, walkTypeof)
	if proc.IsNil() {
		proc.Set(newWalk(cfg))
	} else {
		old := proc.Data.(*walk)
		old.Close()
		proc.Set(newWalk(cfg))
	}

	L.Push(proc)
	return 1
}

func newLuaFileGlob(L *lua.LState) int {
	L.Push(newFileGlob(L))
	return 1
}

func WithEnv(env assert.Environment) {
	xEnv = env
	file := lua.NewUserKV()
	file.Set("open", lua.NewFunction(newLuaFileOpen))
	file.Set("dir", lua.NewFunction(newLuaFileDir))
	file.Set("stat", lua.NewFunction(newLuaFileStat))
	file.Set("walk", lua.NewFunction(newLuaFileWalk))
	file.Set("glob", lua.NewFunction(newLuaFileGlob))
	env.Global("file", file)
}
