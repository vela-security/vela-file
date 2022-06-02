package file

import (
	"fmt"
	"github.com/vela-security/vela-public/lua"
	"github.com/vela-security/vela-public/pipe"
)

var subscript int = 0

type walkConfig struct {
	dir    bool
	deep   int
	name   string
	path   []string
	ignore []func(string) bool
	filter []func(string) bool
	limit  int
	co     *lua.LState
	pipe   *pipe.Px
}

//local w = file.walk("/var/logs")
//w.ignore("*.log")
//w.filter("*java*")

func (cfg *walkConfig) append(path string) {
	cfg.path = append(cfg.path, path)
}

func newWalkConfig(L *lua.LState) *walkConfig {
	subscript++

	n := L.GetTop()
	if n == 0 {
		L.RaiseError("not found path")
		return nil
	}

	w := &walkConfig{
		co:    xEnv.Clone(L),
		dir:   false,
		name:  fmt.Sprintf("walk.%d", subscript),
		pipe:  pipe.New(pipe.Env(xEnv)),
		limit: 0,
	}

	for i := 1; i <= n; i++ {
		w.append(L.CheckString(i))
	}

	return w
}
